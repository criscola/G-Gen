# Copyright (c) 2018 Ultimaker B.V.
# Cura is released under the terms of the LGPLv3 or higher.

from UM.FileHandler.FileWriter import FileWriter #To choose based on the output file mode (text vs. binary).
from UM.FileHandler.WriteFileJob import WriteFileJob #To call the file writer asynchronously.
from UM.Logger import Logger
from UM.Application import Application
from UM.Settings.ContainerRegistry import ContainerRegistry
from UM.i18n import i18nCatalog
from UM.Message import Message
from UM.Qt.Duration import Duration, DurationFormat
from UM.OutputDevice import OutputDeviceError #To show that something went wrong when writing.
from UM.Scene.SceneNode import SceneNode #For typing.
from UM.Version import Version #To check against firmware versions for support.

from cura.PrinterOutput.NetworkedPrinterOutputDevice import NetworkedPrinterOutputDevice, AuthState
from cura.PrinterOutput.PrinterOutputModel import PrinterOutputModel
from cura.PrinterOutput.PrintJobOutputModel import PrintJobOutputModel
from cura.PrinterOutput.MaterialOutputModel import MaterialOutputModel
from cura.PrinterOutput.NetworkCamera import NetworkCamera

from .ClusterUM3PrinterOutputController import ClusterUM3PrinterOutputController

from PyQt5.QtNetwork import QNetworkRequest, QNetworkReply
from PyQt5.QtGui import QDesktopServices
from PyQt5.QtCore import pyqtSlot, QUrl, pyqtSignal, pyqtProperty, QObject

from time import time, sleep
from datetime import datetime
from typing import Optional, Dict, List

import io #To create the correct buffers for sending data to the printer.
import json
import os

i18n_catalog = i18nCatalog("cura")


class ClusterUM3OutputDevice(NetworkedPrinterOutputDevice):
    printJobsChanged = pyqtSignal()
    activePrinterChanged = pyqtSignal()

    # This is a bit of a hack, as the notify can only use signals that are defined by the class that they are in.
    # Inheritance doesn't seem to work. Tying them together does work, but i'm open for better suggestions.
    clusterPrintersChanged = pyqtSignal()

    def __init__(self, device_id, address, properties, parent = None):
        super().__init__(device_id = device_id, address = address, properties=properties, parent = parent)
        self._api_prefix = "/cluster-api/v1/"

        self._number_of_extruders = 2

        self._print_jobs = []

        self._monitor_view_qml_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "ClusterMonitorItem.qml")
        self._control_view_qml_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "ClusterControlItem.qml")

        # See comments about this hack with the clusterPrintersChanged signal
        self.printersChanged.connect(self.clusterPrintersChanged)

        self._accepts_commands = True

        # Cluster does not have authentication, so default to authenticated
        self._authentication_state = AuthState.Authenticated

        self._error_message = None
        self._progress_message = None

        self._active_printer = None  # type: Optional[PrinterOutputModel]

        self._printer_selection_dialog = None

        self.setPriority(3)  # Make sure the output device gets selected above local file output
        self.setName(self._id)
        self.setShortDescription(i18n_catalog.i18nc("@action:button Preceded by 'Ready to'.", "Print over network"))
        self.setDescription(i18n_catalog.i18nc("@properties:tooltip", "Print over network"))

        self.setConnectionText(i18n_catalog.i18nc("@info:status", "Connected over the network"))

        self._printer_uuid_to_unique_name_mapping = {}

        self._finished_jobs = []

        self._cluster_size = int(properties.get(b"cluster_size", 0))

        self._latest_reply_handler = None

    def requestWrite(self, nodes: List[SceneNode], file_name=None, filter_by_machine=False, file_handler=None, **kwargs):
        self.writeStarted.emit(self)

        #Formats supported by this application (file types that we can actually write).
        if file_handler:
            file_formats = file_handler.getSupportedFileTypesWrite()
        else:
            file_formats = Application.getInstance().getMeshFileHandler().getSupportedFileTypesWrite()

        #Create a list from the supported file formats string.
        machine_file_formats = Application.getInstance().getGlobalContainerStack().getMetaDataEntry("file_formats").split(";")
        machine_file_formats = [file_type.strip() for file_type in machine_file_formats]
        #Exception for UM3 firmware version >=4.4: UFP is now supported and should be the preferred file format.
        if "application/x-ufp" not in machine_file_formats and self.printerType == "ultimaker3" and Version(self.firmwareVersion) >= Version("4.4"):
            machine_file_formats = ["application/x-ufp"] + machine_file_formats

        # Take the intersection between file_formats and machine_file_formats.
        format_by_mimetype = {format["mime_type"]: format for format in file_formats}
        file_formats = [format_by_mimetype[mimetype] for mimetype in machine_file_formats] #Keep them ordered according to the preference in machine_file_formats.

        if len(file_formats) == 0:
            Logger.log("e", "There are no file formats available to write with!")
            raise OutputDeviceError.WriteRequestFailedError(i18n_catalog.i18nc("@info:status", "There are no file formats available to write with!"))
        preferred_format = file_formats[0]

        #Just take the first file format available.
        if file_handler is not None:
            writer = file_handler.getWriterByMimeType(preferred_format["mime_type"])
        else:
            writer = Application.getInstance().getMeshFileHandler().getWriterByMimeType(preferred_format["mime_type"])

        #This function pauses with the yield, waiting on instructions on which printer it needs to print with.
        self._sending_job = self._sendPrintJob(writer, preferred_format, nodes)
        self._sending_job.send(None) #Start the generator.

        if len(self._printers) > 1: #We need to ask the user.
            self._spawnPrinterSelectionDialog()
            is_job_sent = True
        else: #Just immediately continue.
            self._sending_job.send("") #No specifically selected printer.
            is_job_sent = self._sending_job.send(None)

    def _spawnPrinterSelectionDialog(self):
        if self._printer_selection_dialog is None:
            path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "PrintWindow.qml")
            self._printer_selection_dialog = Application.getInstance().createQmlComponent(path, {"OutputDevice": self})
        if self._printer_selection_dialog is not None:
            self._printer_selection_dialog.show()

    @pyqtProperty(int, constant=True)
    def clusterSize(self):
        return self._cluster_size

    ##  Allows the user to choose a printer to print with from the printer
    #   selection dialogue.
    #   \param target_printer The name of the printer to target.
    @pyqtSlot(str)
    def selectPrinter(self, target_printer: str = "") -> None:
        self._sending_job.send(target_printer)

    ##  Greenlet to send a job to the printer over the network.
    #
    #   This greenlet gets called asynchronously in requestWrite. It is a
    #   greenlet in order to optionally wait for selectPrinter() to select a
    #   printer.
    #   The greenlet yields exactly three times: First time None,
    #   \param writer The file writer to use to create the data.
    #   \param preferred_format A dictionary containing some information about
    #   what format to write to. This is necessary to create the correct buffer
    #   types and file extension and such.
    def _sendPrintJob(self, writer: FileWriter, preferred_format: Dict, nodes: List[SceneNode]):
        Logger.log("i", "Sending print job to printer.")
        if self._sending_gcode:
            self._error_message = Message(
                i18n_catalog.i18nc("@info:status",
                                   "Sending new jobs (temporarily) blocked, still sending the previous print job."))
            self._error_message.show()
            yield #Wait on the user to select a target printer.
            yield #Wait for the write job to be finished.
            yield False #Return whether this was a success or not.
            yield #Prevent StopIteration.

        self._sending_gcode = True

        target_printer = yield #Potentially wait on the user to select a target printer.

        # Using buffering greatly reduces the write time for many lines of gcode
        if preferred_format["mode"] == FileWriter.OutputMode.TextMode:
            stream = io.StringIO()
        else: #Binary mode.
            stream = io.BytesIO()

        job = WriteFileJob(writer, stream, nodes, preferred_format["mode"])

        self._progress_message = Message(i18n_catalog.i18nc("@info:status", "Sending data to printer"), lifetime = 0, dismissable = False, progress = -1,
                                         title = i18n_catalog.i18nc("@info:title", "Sending Data"))
        self._progress_message.addAction("Abort", i18n_catalog.i18nc("@action:button", "Cancel"), icon = None, description = "")
        self._progress_message.actionTriggered.connect(self._progressMessageActionTriggered)
        self._progress_message.show()

        job.start()

        parts = []

        # If a specific printer was selected, it should be printed with that machine.
        if target_printer:
            target_printer = self._printer_uuid_to_unique_name_mapping[target_printer]
            parts.append(self._createFormPart("name=require_printer_name", bytes(target_printer, "utf-8"), "text/plain"))

        # Add user name to the print_job
        parts.append(self._createFormPart("name=owner", bytes(self._getUserName(), "utf-8"), "text/plain"))

        file_name = Application.getInstance().getPrintInformation().jobName + "." + preferred_format["extension"]

        while not job.isFinished():
            sleep(0.1)
        output = stream.getvalue() #Either str or bytes depending on the output mode.
        if isinstance(stream, io.StringIO):
            output = output.encode("utf-8")

        parts.append(self._createFormPart("name=\"file\"; filename=\"%s\"" % file_name, output))

        self._latest_reply_handler = self.postFormWithParts("print_jobs/", parts, onFinished=self._onPostPrintJobFinished, onProgress=self._onUploadPrintJobProgress)

        yield True #Return that we had success!
        yield #To prevent having to catch the StopIteration exception.

    @pyqtProperty(QObject, notify=activePrinterChanged)
    def activePrinter(self) -> Optional[PrinterOutputModel]:
        return self._active_printer

    @pyqtSlot(QObject)
    def setActivePrinter(self, printer: Optional[PrinterOutputModel]):
        if self._active_printer != printer:
            if self._active_printer and self._active_printer.camera:
                self._active_printer.camera.stop()
            self._active_printer = printer
            self.activePrinterChanged.emit()

    def _onPostPrintJobFinished(self, reply):
        self._progress_message.hide()
        self._compressing_gcode = False
        self._sending_gcode = False

    def _onUploadPrintJobProgress(self, bytes_sent:int, bytes_total:int):
        if bytes_total > 0:
            new_progress = bytes_sent / bytes_total * 100
            # Treat upload progress as response. Uploading can take more than 10 seconds, so if we don't, we can get
            # timeout responses if this happens.
            self._last_response_time = time()
            if new_progress > self._progress_message.getProgress():
                self._progress_message.show()  # Ensure that the message is visible.
                self._progress_message.setProgress(bytes_sent / bytes_total * 100)

            # If successfully sent:
            if bytes_sent == bytes_total:
                # Show a confirmation to the user so they know the job was sucessful and provide the option to switch to the
                # monitor tab.
                self._success_message = Message(
                    i18n_catalog.i18nc("@info:status", "Print job was successfully sent to the printer."),
                    lifetime=5, dismissable=True,
                    title=i18n_catalog.i18nc("@info:title", "Data Sent"))
                self._success_message.addAction("View", i18n_catalog.i18nc("@action:button", "View in Monitor"), icon=None,
                                                description="")
                self._success_message.actionTriggered.connect(self._successMessageActionTriggered)
                self._success_message.show()
        else:
            self._progress_message.setProgress(0)
            self._progress_message.hide()

    def _progressMessageActionTriggered(self, message_id: Optional[str]=None, action_id: Optional[str]=None) -> None:
        if action_id == "Abort":
            Logger.log("d", "User aborted sending print to remote.")
            self._progress_message.hide()
            self._compressing_gcode = False
            self._sending_gcode = False
            Application.getInstance().getController().setActiveStage("PrepareStage")

            # After compressing the sliced model Cura sends data to printer, to stop receiving updates from the request
            # the "reply" should be disconnected
            if self._latest_reply_handler:
                self._latest_reply_handler.disconnect()
                self._latest_reply_handler = None

    def _successMessageActionTriggered(self, message_id: Optional[str]=None, action_id: Optional[str]=None) -> None:
        if action_id == "View":
            Application.getInstance().getController().setActiveStage("MonitorStage")

    @pyqtSlot()
    def openPrintJobControlPanel(self) -> None:
        Logger.log("d", "Opening print job control panel...")
        QDesktopServices.openUrl(QUrl("http://" + self._address + "/print_jobs"))

    @pyqtSlot()
    def openPrinterControlPanel(self) -> None:
        Logger.log("d", "Opening printer control panel...")
        QDesktopServices.openUrl(QUrl("http://" + self._address + "/printers"))

    @pyqtProperty("QVariantList", notify=printJobsChanged)
    def printJobs(self)-> List[PrintJobOutputModel] :
        return self._print_jobs

    @pyqtProperty("QVariantList", notify=printJobsChanged)
    def queuedPrintJobs(self) -> List[PrintJobOutputModel]:
        return [print_job for print_job in self._print_jobs if print_job.state == "queued"]

    @pyqtProperty("QVariantList", notify=printJobsChanged)
    def activePrintJobs(self) -> List[PrintJobOutputModel]:
        return [print_job for print_job in self._print_jobs if print_job.assignedPrinter is not None and print_job.state != "queued"]

    @pyqtProperty("QVariantList", notify=clusterPrintersChanged)
    def connectedPrintersTypeCount(self) -> List[PrinterOutputModel]:
        printer_count = {}
        for printer in self._printers:
            if printer.type in printer_count:
                printer_count[printer.type] += 1
            else:
                printer_count[printer.type] = 1
        result = []
        for machine_type in printer_count:
            result.append({"machine_type": machine_type, "count": printer_count[machine_type]})
        return result

    @pyqtSlot(int, result=str)
    def formatDuration(self, seconds: int) -> str:
        return Duration(seconds).getDisplayString(DurationFormat.Format.Short)

    @pyqtSlot(int, result=str)
    def getTimeCompleted(self, time_remaining: int) -> str:
        current_time = time()
        datetime_completed = datetime.fromtimestamp(current_time + time_remaining)
        return "{hour:02d}:{minute:02d}".format(hour=datetime_completed.hour, minute=datetime_completed.minute)

    @pyqtSlot(int, result=str)
    def getDateCompleted(self, time_remaining: int) -> str:
        current_time = time()
        datetime_completed = datetime.fromtimestamp(current_time + time_remaining)
        return (datetime_completed.strftime("%a %b ") + "{day}".format(day=datetime_completed.day)).upper()

    def _printJobStateChanged(self) -> None:
        username = self._getUserName()

        if username is None:
            return  # We only want to show notifications if username is set.

        finished_jobs = [job for job in self._print_jobs if job.state == "wait_cleanup"]

        newly_finished_jobs = [job for job in finished_jobs if job not in self._finished_jobs and job.owner == username]
        for job in newly_finished_jobs:
            if job.assignedPrinter:
                job_completed_text = i18n_catalog.i18nc("@info:status", "Printer '{printer_name}' has finished printing '{job_name}'.".format(printer_name=job.assignedPrinter.name, job_name = job.name))
            else:
                job_completed_text =  i18n_catalog.i18nc("@info:status", "The print job '{job_name}' was finished.".format(job_name = job.name))
            job_completed_message = Message(text=job_completed_text, title = i18n_catalog.i18nc("@info:status", "Print finished"))
            job_completed_message.show()

        # Ensure UI gets updated
        self.printJobsChanged.emit()

        # Keep a list of all completed jobs so we know if something changed next time.
        self._finished_jobs = finished_jobs

    def _update(self) -> None:
        if not super()._update():
            return
        self.get("printers/", onFinished=self._onGetPrintersDataFinished)
        self.get("print_jobs/", onFinished=self._onGetPrintJobsFinished)

    def _onGetPrintJobsFinished(self, reply: QNetworkReply) -> None:
        if not checkValidGetReply(reply):
            return

        result = loadJsonFromReply(reply)
        if result is None:
            return

        print_jobs_seen = []
        job_list_changed = False
        for print_job_data in result:
            print_job = findByKey(self._print_jobs, print_job_data["uuid"])

            if print_job is None:
                print_job = self._createPrintJobModel(print_job_data)
                job_list_changed = True

            self._updatePrintJob(print_job, print_job_data)

            if print_job.state != "queued":  # Print job should be assigned to a printer.
                if print_job.state in ["failed", "finished", "aborted"]:
                    # Print job was already completed, so don't attach it to a printer.
                    printer = None
                else:
                    printer = self._getPrinterByKey(print_job_data["printer_uuid"])
            else:  # The job can "reserve" a printer if some changes are required.
                printer = self._getPrinterByKey(print_job_data["assigned_to"])

            if printer:
                printer.updateActivePrintJob(print_job)

            print_jobs_seen.append(print_job)

        # Check what jobs need to be removed.
        removed_jobs = [print_job for print_job in self._print_jobs if print_job not in print_jobs_seen]

        for removed_job in removed_jobs:
            job_list_changed |= self._removeJob(removed_job)

        if job_list_changed:
            self.printJobsChanged.emit()  # Do a single emit for all print job changes.

    def _onGetPrintersDataFinished(self, reply: QNetworkReply) -> None:
        if not checkValidGetReply(reply):
            return

        result = loadJsonFromReply(reply)
        if result is None:
            return

        printer_list_changed = False
        printers_seen = []

        for printer_data in result:
            printer = findByKey(self._printers, printer_data["uuid"])

            if printer is None:
                printer = self._createPrinterModel(printer_data)
                printer_list_changed = True

            printers_seen.append(printer)

            self._updatePrinter(printer, printer_data)

        removed_printers = [printer for printer in self._printers if printer not in printers_seen]
        for printer in removed_printers:
            self._removePrinter(printer)

        if removed_printers or printer_list_changed:
            self.printersChanged.emit()

    def _createPrinterModel(self, data: Dict) -> PrinterOutputModel:
        printer = PrinterOutputModel(output_controller=ClusterUM3PrinterOutputController(self),
                                     number_of_extruders=self._number_of_extruders)
        printer.setCamera(NetworkCamera("http://" + data["ip_address"] + ":8080/?action=stream"))
        self._printers.append(printer)
        return printer

    def _createPrintJobModel(self, data: Dict) -> PrintJobOutputModel:
        print_job = PrintJobOutputModel(output_controller=ClusterUM3PrinterOutputController(self),
                                        key=data["uuid"], name= data["name"])
        print_job.stateChanged.connect(self._printJobStateChanged)
        self._print_jobs.append(print_job)
        return print_job

    def _updatePrintJob(self, print_job: PrintJobOutputModel, data: Dict) -> None:
        print_job.updateTimeTotal(data["time_total"])
        print_job.updateTimeElapsed(data["time_elapsed"])
        print_job.updateState(data["status"])
        print_job.updateOwner(data["owner"])

    def _updatePrinter(self, printer: PrinterOutputModel, data: Dict) -> None:
        # For some unknown reason the cluster wants UUID for everything, except for sending a job directly to a printer.
        # Then we suddenly need the unique name. So in order to not have to mess up all the other code, we save a mapping.
        self._printer_uuid_to_unique_name_mapping[data["uuid"]] = data["unique_name"]

        definitions = ContainerRegistry.getInstance().findDefinitionContainers(name = data["machine_variant"])
        if not definitions:
            Logger.log("w", "Unable to find definition for machine variant %s", data["machine_variant"])
            return

        machine_definition = definitions[0]

        printer.updateName(data["friendly_name"])
        printer.updateKey(data["uuid"])
        printer.updateType(data["machine_variant"])

        # Do not store the buildplate information that comes from connect if the current printer has not buildplate information
        if "build_plate" in data and machine_definition.getMetaDataEntry("has_variant_buildplates", False):
            printer.updateBuildplateName(data["build_plate"]["type"])
        if not data["enabled"]:
            printer.updateState("disabled")
        else:
            printer.updateState(data["status"])

        for index in range(0, self._number_of_extruders):
            extruder = printer.extruders[index]
            try:
                extruder_data = data["configuration"][index]
            except IndexError:
                break

            extruder.updateHotendID(extruder_data.get("print_core_id", ""))

            material_data = extruder_data["material"]
            if extruder.activeMaterial is None or extruder.activeMaterial.guid != material_data["guid"]:
                containers = ContainerRegistry.getInstance().findInstanceContainers(type="material",
                                                                                    GUID=material_data["guid"])
                if containers:
                    color = containers[0].getMetaDataEntry("color_code")
                    brand = containers[0].getMetaDataEntry("brand")
                    material_type = containers[0].getMetaDataEntry("material")
                    name = containers[0].getName()
                else:
                    Logger.log("w",
                               "Unable to find material with guid {guid}. Using data as provided by cluster".format(
                                   guid=material_data["guid"]))
                    color = material_data["color"]
                    brand = material_data["brand"]
                    material_type = material_data["material"]
                    name = "Empty" if material_data["material"] == "empty" else "Unknown"

                material = MaterialOutputModel(guid=material_data["guid"], type=material_type,
                                               brand=brand, color=color, name=name)
                extruder.updateActiveMaterial(material)

    def _removeJob(self, job: PrintJobOutputModel):
        if job not in self._print_jobs:
            return False

        if job.assignedPrinter:
            job.assignedPrinter.updateActivePrintJob(None)
            job.stateChanged.disconnect(self._printJobStateChanged)
        self._print_jobs.remove(job)

        return True

    def _removePrinter(self, printer: PrinterOutputModel):
        self._printers.remove(printer)
        if self._active_printer == printer:
            self._active_printer = None
            self.activePrinterChanged.emit()


def loadJsonFromReply(reply):
    try:
        result = json.loads(bytes(reply.readAll()).decode("utf-8"))
    except json.decoder.JSONDecodeError:
        Logger.logException("w", "Unable to decode JSON from reply.")
        return
    return result


def checkValidGetReply(reply):
    status_code = reply.attribute(QNetworkRequest.HttpStatusCodeAttribute)

    if status_code != 200:
        Logger.log("w", "Got status code {status_code} while trying to get data".format(status_code=status_code))
        return False
    return True


def findByKey(list, key):
    for item in list:
        if item.key == key:
            return item