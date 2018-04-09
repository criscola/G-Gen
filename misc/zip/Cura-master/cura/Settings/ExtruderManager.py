# Copyright (c) 2017 Ultimaker B.V.
# Cura is released under the terms of the LGPLv3 or higher.

from PyQt5.QtCore import pyqtSignal, pyqtProperty, QObject, QVariant  # For communicating data and events to Qt.
from UM.FlameProfiler import pyqtSlot

from UM.Application import Application  # To get the global container stack to find the current machine.
from UM.Logger import Logger
from UM.Scene.Iterator.DepthFirstIterator import DepthFirstIterator
from UM.Scene.SceneNode import SceneNode
from UM.Scene.Selection import Selection
from UM.Scene.Iterator.BreadthFirstIterator import BreadthFirstIterator
from UM.Settings.ContainerRegistry import ContainerRegistry  # Finding containers by ID.
from UM.Settings.SettingFunction import SettingFunction
from UM.Settings.SettingInstance import SettingInstance
from UM.Settings.ContainerStack import ContainerStack
from UM.Settings.PropertyEvaluationContext import PropertyEvaluationContext
from typing import Optional, List, TYPE_CHECKING, Union

if TYPE_CHECKING:
    from cura.Settings.ExtruderStack import ExtruderStack
    from cura.Settings.GlobalStack import GlobalStack


##  Manages all existing extruder stacks.
#
#   This keeps a list of extruder stacks for each machine.
class ExtruderManager(QObject):

    ##  Registers listeners and such to listen to changes to the extruders.
    def __init__(self, parent = None):
        super().__init__(parent)

        self._application = Application.getInstance()

        self._extruder_trains = {}  # Per machine, a dictionary of extruder container stack IDs. Only for separately defined extruders.
        self._active_extruder_index = -1  # Indicates the index of the active extruder stack. -1 means no active extruder stack
        self._selected_object_extruders = []
        self._addCurrentMachineExtruders()

        #Application.getInstance().globalContainerStackChanged.connect(self._globalContainerStackChanged)
        Selection.selectionChanged.connect(self.resetSelectedObjectExtruders)

    ##  Signal to notify other components when the list of extruders for a machine definition changes.
    extrudersChanged = pyqtSignal(QVariant)

    ##  Notify when the user switches the currently active extruder.
    activeExtruderChanged = pyqtSignal()

    ##  Gets the unique identifier of the currently active extruder stack.
    #
    #   The currently active extruder stack is the stack that is currently being
    #   edited.
    #
    #   \return The unique ID of the currently active extruder stack.
    @pyqtProperty(str, notify = activeExtruderChanged)
    def activeExtruderStackId(self) -> Optional[str]:
        if not Application.getInstance().getGlobalContainerStack():
            return None  # No active machine, so no active extruder.
        try:
            return self._extruder_trains[Application.getInstance().getGlobalContainerStack().getId()][str(self._active_extruder_index)].getId()
        except KeyError:  # Extruder index could be -1 if the global tab is selected, or the entry doesn't exist if the machine definition is wrong.
            return None

    ##  Return extruder count according to extruder trains.
    @pyqtProperty(int, notify = extrudersChanged)
    def extruderCount(self):
        if not Application.getInstance().getGlobalContainerStack():
            return 0  # No active machine, so no extruders.
        try:
            return len(self._extruder_trains[Application.getInstance().getGlobalContainerStack().getId()])
        except KeyError:
            return 0

    ##  Gets a dict with the extruder stack ids with the extruder number as the key.
    @pyqtProperty("QVariantMap", notify = extrudersChanged)
    def extruderIds(self):
        extruder_stack_ids = {}

        global_stack_id = Application.getInstance().getGlobalContainerStack().getId()

        if global_stack_id in self._extruder_trains:
            for position in self._extruder_trains[global_stack_id]:
                extruder_stack_ids[position] = self._extruder_trains[global_stack_id][position].getId()

        return extruder_stack_ids

    @pyqtSlot(str, result = str)
    def getQualityChangesIdByExtruderStackId(self, extruder_stack_id: str) -> str:
        for position in self._extruder_trains[Application.getInstance().getGlobalContainerStack().getId()]:
            extruder = self._extruder_trains[Application.getInstance().getGlobalContainerStack().getId()][position]
            if extruder.getId() == extruder_stack_id:
                return extruder.qualityChanges.getId()

    ##  The instance of the singleton pattern.
    #
    #   It's None if the extruder manager hasn't been created yet.
    __instance = None

    @staticmethod
    def createExtruderManager():
        return ExtruderManager().getInstance()

    ##  Gets an instance of the extruder manager, or creates one if no instance
    #   exists yet.
    #
    #   This is an implementation of singleton. If an extruder manager already
    #   exists, it is re-used.
    #
    #   \return The extruder manager.
    @classmethod
    def getInstance(cls) -> "ExtruderManager":
        if not cls.__instance:
            cls.__instance = ExtruderManager()
        return cls.__instance

    ##  Changes the active extruder by index.
    #
    #   \param index The index of the new active extruder.
    @pyqtSlot(int)
    def setActiveExtruderIndex(self, index: int) -> None:
        self._active_extruder_index = index
        self.activeExtruderChanged.emit()

    @pyqtProperty(int, notify = activeExtruderChanged)
    def activeExtruderIndex(self) -> int:
        return self._active_extruder_index

    ##  Gets the extruder name of an extruder of the currently active machine.
    #
    #   \param index The index of the extruder whose name to get.
    @pyqtSlot(int, result = str)
    def getExtruderName(self, index):
        try:
            return list(self.getActiveExtruderStacks())[index].getName()
        except IndexError:
            return ""

    ## Emitted whenever the selectedObjectExtruders property changes.
    selectedObjectExtrudersChanged = pyqtSignal()

    ##  Provides a list of extruder IDs used by the current selected objects.
    @pyqtProperty("QVariantList", notify = selectedObjectExtrudersChanged)
    def selectedObjectExtruders(self) -> List[str]:
        if not self._selected_object_extruders:
            object_extruders = set()

            # First, build a list of the actual selected objects (including children of groups, excluding group nodes)
            selected_nodes = []
            for node in Selection.getAllSelectedObjects():
                if node.callDecoration("isGroup"):
                    for grouped_node in BreadthFirstIterator(node):
                        if grouped_node.callDecoration("isGroup"):
                            continue

                        selected_nodes.append(grouped_node)
                else:
                    selected_nodes.append(node)

            # Then, figure out which nodes are used by those selected nodes.
            global_stack = Application.getInstance().getGlobalContainerStack()
            current_extruder_trains = self._extruder_trains.get(global_stack.getId())
            for node in selected_nodes:
                extruder = node.callDecoration("getActiveExtruder")
                if extruder:
                    object_extruders.add(extruder)
                elif current_extruder_trains:
                    object_extruders.add(current_extruder_trains["0"].getId())

            self._selected_object_extruders = list(object_extruders)

        return self._selected_object_extruders

    ##  Reset the internal list used for the selectedObjectExtruders property
    #
    #   This will trigger a recalculation of the extruders used for the
    #   selection.
    def resetSelectedObjectExtruders(self) -> None:
        self._selected_object_extruders = []
        self.selectedObjectExtrudersChanged.emit()

    @pyqtSlot(result = QObject)
    def getActiveExtruderStack(self) -> Optional["ExtruderStack"]:
        global_container_stack = Application.getInstance().getGlobalContainerStack()

        if global_container_stack:
            if global_container_stack.getId() in self._extruder_trains:
                if str(self._active_extruder_index) in self._extruder_trains[global_container_stack.getId()]:
                    return self._extruder_trains[global_container_stack.getId()][str(self._active_extruder_index)]

        return None

    ##  Get an extruder stack by index
    def getExtruderStack(self, index) -> Optional["ExtruderStack"]:
        global_container_stack = Application.getInstance().getGlobalContainerStack()
        if global_container_stack:
            if global_container_stack.getId() in self._extruder_trains:
                if str(index) in self._extruder_trains[global_container_stack.getId()]:
                    return self._extruder_trains[global_container_stack.getId()][str(index)]
        return None

    ##  Get all extruder stacks
    def getExtruderStacks(self) -> List["ExtruderStack"]:
        result = []
        for i in range(self.extruderCount):
            result.append(self.getExtruderStack(i))
        return result

    def registerExtruder(self, extruder_train, machine_id):
        changed = False

        if machine_id not in self._extruder_trains:
            self._extruder_trains[machine_id] = {}
            changed = True

        # do not register if an extruder has already been registered at the position on this machine
        if any(item.getId() == extruder_train.getId() for item in self._extruder_trains[machine_id].values()):
            Logger.log("w", "Extruder [%s] has already been registered on machine [%s], not doing anything",
                       extruder_train.getId(), machine_id)
            return

        if extruder_train:
            self._extruder_trains[machine_id][extruder_train.getMetaDataEntry("position")] = extruder_train
            changed = True
        if changed:
            self.extrudersChanged.emit(machine_id)

    def getAllExtruderValues(self, setting_key):
        return self.getAllExtruderSettings(setting_key, "value")

    ##  Gets a property of a setting for all extruders.
    #
    #   \param setting_key  \type{str} The setting to get the property of.
    #   \param property  \type{str} The property to get.
    #   \return \type{List} the list of results
    def getAllExtruderSettings(self, setting_key: str, prop: str):
        result = []
        for index in self.extruderIds:
            extruder_stack_id = self.extruderIds[str(index)]
            extruder_stack = ContainerRegistry.getInstance().findContainerStacks(id = extruder_stack_id)[0]
            result.append(extruder_stack.getProperty(setting_key, prop))
        return result

    def extruderValueWithDefault(self, value):
        machine_manager = self._application.getMachineManager()
        if value == "-1":
            return machine_manager.defaultExtruderPosition
        else:
            return value

    ##  Gets the extruder stacks that are actually being used at the moment.
    #
    #   An extruder stack is being used if it is the extruder to print any mesh
    #   with, or if it is the support infill extruder, the support interface
    #   extruder, or the bed adhesion extruder.
    #
    #   If there are no extruders, this returns the global stack as a singleton
    #   list.
    #
    #   \return A list of extruder stacks.
    def getUsedExtruderStacks(self) -> List["ContainerStack"]:
        global_stack = self._application.getGlobalContainerStack()
        container_registry = ContainerRegistry.getInstance()

        used_extruder_stack_ids = set()

        # Get the extruders of all meshes in the scene
        support_enabled = False
        support_bottom_enabled = False
        support_roof_enabled = False

        scene_root = Application.getInstance().getController().getScene().getRoot()

        # If no extruders are registered in the extruder manager yet, return an empty array
        if len(self.extruderIds) == 0:
            return []

        # Get the extruders of all printable meshes in the scene
        meshes = [node for node in DepthFirstIterator(scene_root) if isinstance(node, SceneNode) and node.isSelectable()]
        for mesh in meshes:
            extruder_stack_id = mesh.callDecoration("getActiveExtruder")
            if not extruder_stack_id:
                # No per-object settings for this node
                extruder_stack_id = self.extruderIds["0"]
            used_extruder_stack_ids.add(extruder_stack_id)

            # Get whether any of them use support.
            stack_to_use = mesh.callDecoration("getStack")  # if there is a per-mesh stack, we use it
            if not stack_to_use:
                # if there is no per-mesh stack, we use the build extruder for this mesh
                stack_to_use = container_registry.findContainerStacks(id = extruder_stack_id)[0]

            support_enabled |= stack_to_use.getProperty("support_enable", "value")
            support_bottom_enabled |= stack_to_use.getProperty("support_bottom_enable", "value")
            support_roof_enabled |= stack_to_use.getProperty("support_roof_enable", "value")

            # Check limit to extruders
            limit_to_extruder_feature_list = ["wall_0_extruder_nr",
                                              "wall_x_extruder_nr",
                                              "roofing_extruder_nr",
                                              "top_bottom_extruder_nr",
                                              "infill_extruder_nr",
                                              ]
            for extruder_nr_feature_name in limit_to_extruder_feature_list:
                extruder_nr = int(global_stack.getProperty(extruder_nr_feature_name, "value"))
                if extruder_nr == -1:
                    continue
                used_extruder_stack_ids.add(self.extruderIds[str(extruder_nr)])

        # Check support extruders
        if support_enabled:
            used_extruder_stack_ids.add(self.extruderIds[self.extruderValueWithDefault(str(global_stack.getProperty("support_infill_extruder_nr", "value")))])
            used_extruder_stack_ids.add(self.extruderIds[self.extruderValueWithDefault(str(global_stack.getProperty("support_extruder_nr_layer_0", "value")))])
            if support_bottom_enabled:
                used_extruder_stack_ids.add(self.extruderIds[self.extruderValueWithDefault(str(global_stack.getProperty("support_bottom_extruder_nr", "value")))])
            if support_roof_enabled:
                used_extruder_stack_ids.add(self.extruderIds[self.extruderValueWithDefault(str(global_stack.getProperty("support_roof_extruder_nr", "value")))])

        # The platform adhesion extruder. Not used if using none.
        if global_stack.getProperty("adhesion_type", "value") != "none":
            extruder_nr = str(global_stack.getProperty("adhesion_extruder_nr", "value"))
            if extruder_nr == "-1":
                extruder_nr = Application.getInstance().getMachineManager().defaultExtruderPosition
            used_extruder_stack_ids.add(self.extruderIds[extruder_nr])

        try:
            return [container_registry.findContainerStacks(id = stack_id)[0] for stack_id in used_extruder_stack_ids]
        except IndexError:  # One or more of the extruders was not found.
            Logger.log("e", "Unable to find one or more of the extruders in %s", used_extruder_stack_ids)
            return []

    ##  Removes the container stack and user profile for the extruders for a specific machine.
    #
    #   \param machine_id The machine to remove the extruders for.
    def removeMachineExtruders(self, machine_id: str):
        for extruder in self.getMachineExtruders(machine_id):
            ContainerRegistry.getInstance().removeContainer(extruder.userChanges.getId())
            ContainerRegistry.getInstance().removeContainer(extruder.getId())
        if machine_id in self._extruder_trains:
            del self._extruder_trains[machine_id]

    ##  Returns extruders for a specific machine.
    #
    #   \param machine_id The machine to get the extruders of.
    def getMachineExtruders(self, machine_id: str):
        if machine_id not in self._extruder_trains:
            return []
        return [self._extruder_trains[machine_id][name] for name in self._extruder_trains[machine_id]]

    ##  Returns a list containing the global stack and active extruder stacks.
    #
    #   The first element is the global container stack, followed by any extruder stacks.
    #   \return \type{List[ContainerStack]}
    def getActiveGlobalAndExtruderStacks(self) -> Optional[List[Union["ExtruderStack", "GlobalStack"]]]:
        global_stack = Application.getInstance().getGlobalContainerStack()
        if not global_stack:
            return None

        result = [global_stack]
        result.extend(self.getActiveExtruderStacks())
        return result

    ##  Returns the list of active extruder stacks, taking into account the machine extruder count.
    #
    #   \return \type{List[ContainerStack]} a list of
    def getActiveExtruderStacks(self) -> List["ExtruderStack"]:
        global_stack = Application.getInstance().getGlobalContainerStack()
        if not global_stack:
            return []

        result = []
        if global_stack.getId() in self._extruder_trains:
            for extruder in sorted(self._extruder_trains[global_stack.getId()]):
                result.append(self._extruder_trains[global_stack.getId()][extruder])

        machine_extruder_count = global_stack.getProperty("machine_extruder_count", "value")

        return result[:machine_extruder_count]

    def _globalContainerStackChanged(self) -> None:
        # If the global container changed, the machine changed and might have extruders that were not registered yet
        self._addCurrentMachineExtruders()

        self.resetSelectedObjectExtruders()

    ##  Adds the extruders of the currently active machine.
    def _addCurrentMachineExtruders(self) -> None:
        global_stack = self._application.getGlobalContainerStack()
        extruders_changed = False

        if global_stack:
            container_registry = ContainerRegistry.getInstance()
            global_stack_id = global_stack.getId()

            # Gets the extruder trains that we just created as well as any that still existed.
            extruder_trains = container_registry.findContainerStacks(type = "extruder_train", machine = global_stack_id)

            # Make sure the extruder trains for the new machine can be placed in the set of sets
            if global_stack_id not in self._extruder_trains:
                self._extruder_trains[global_stack_id] = {}
                extruders_changed = True

            # Register the extruder trains by position
            for extruder_train in extruder_trains:
                self._extruder_trains[global_stack_id][extruder_train.getMetaDataEntry("position")] = extruder_train

                # regardless of what the next stack is, we have to set it again, because of signal routing. ???
                extruder_train.setParent(global_stack)
                extruder_train.setNextStack(global_stack)
                extruders_changed = True

            self._fixMaterialDiameterAndNozzleSize(global_stack, extruder_trains)
            if extruders_changed:
                self.extrudersChanged.emit(global_stack_id)
                self.setActiveExtruderIndex(0)

    #
    # This function tries to fix the problem with per-extruder-settable nozzle size and material diameter problems
    # in early versions (3.0 - 3.2.1).
    #
    # In earlier versions, "nozzle size" and "material diameter" are only applicable to the complete machine, so all
    # extruders share the same values. In this case, "nozzle size" and "material diameter" are saved in the
    # GlobalStack's DefinitionChanges container.
    #
    # Later, we could have different "nozzle size" for each extruder, but "material diameter" could only be set for
    # the entire machine. In this case, "nozzle size" should be saved in each ExtruderStack's DefinitionChanges, but
    # "material diameter" still remains in the GlobalStack's DefinitionChanges.
    #
    # Lateer, both "nozzle size" and "material diameter" are settable per-extruder, and both settings should be saved
    # in the ExtruderStack's DefinitionChanges.
    #
    # There were some bugs in upgrade so the data weren't saved correct as described above. This function tries fix
    # this.
    #
    # One more thing is about material diameter and single-extrusion machines. Most single-extrusion machines don't
    # specifically define their extruder definition, so they reuse "fdmextruder", but for those machines, they may
    # define "material diameter = 1.75" in their machine definition, but in "fdmextruder", it's still "2.85". This
    # causes a problem with incorrect default values.
    #
    # This is also fixed here in this way: If no "material diameter" is specified, it will look for the default value
    # in both the Extruder's definition and the Global's definition. If 2 values don't match, we will use the value
    # from the Global definition by setting it in the Extruder's DefinitionChanges container.
    #
    def _fixMaterialDiameterAndNozzleSize(self, global_stack, extruder_stack_list):
        keys_to_copy = ["material_diameter", "machine_nozzle_size"]  # these will be copied over to all extruders

        extruder_positions_to_update = set()
        for extruder_stack in extruder_stack_list:
            for key in keys_to_copy:
                # Only copy the value when this extruder doesn't have the value.
                if extruder_stack.definitionChanges.hasProperty(key, "value"):
                    continue

                setting_value_in_global_def_changes = global_stack.definitionChanges.getProperty(key, "value")
                setting_value_in_global_def = global_stack.definition.getProperty(key, "value")
                setting_value = setting_value_in_global_def
                if setting_value_in_global_def_changes is not None:
                    setting_value = setting_value_in_global_def_changes
                if setting_value == extruder_stack.definition.getProperty(key, "value"):
                    continue

                setting_definition = global_stack.getSettingDefinition(key)
                new_instance = SettingInstance(setting_definition, extruder_stack.definitionChanges)
                new_instance.setProperty("value", setting_value)
                new_instance.resetState()  # Ensure that the state is not seen as a user state.
                extruder_stack.definitionChanges.addInstance(new_instance)
                extruder_stack.definitionChanges.setDirty(True)

                # Make sure the material diameter is up to date for the extruder stack.
                if key == "material_diameter":
                    position = int(extruder_stack.getMetaDataEntry("position"))
                    extruder_positions_to_update.add(position)

        # We have to remove those settings here because we know that those values have been copied to all
        # the extruders at this point.
        for key in keys_to_copy:
            if global_stack.definitionChanges.hasProperty(key, "value"):
                global_stack.definitionChanges.removeInstance(key, postpone_emit = True)

        # Update material diameter for extruders
        for position in extruder_positions_to_update:
            self.updateMaterialForDiameter(position, global_stack = global_stack)

    ##  Get all extruder values for a certain setting.
    #
    #   This is exposed to SettingFunction so it can be used in value functions.
    #
    #   \param key The key of the setting to retrieve values for.
    #
    #   \return A list of values for all extruders. If an extruder does not have a value, it will not be in the list.
    #           If no extruder has the value, the list will contain the global value.
    @staticmethod
    def getExtruderValues(key):
        global_stack = Application.getInstance().getGlobalContainerStack()

        result = []
        for extruder in ExtruderManager.getInstance().getMachineExtruders(global_stack.getId()):
            if not extruder.isEnabled:
                continue
            # only include values from extruders that are "active" for the current machine instance
            if int(extruder.getMetaDataEntry("position")) >= global_stack.getProperty("machine_extruder_count", "value"):
                continue

            value = extruder.getRawProperty(key, "value")

            if value is None:
                continue

            if isinstance(value, SettingFunction):
                value = value(extruder)

            result.append(value)

        if not result:
            result.append(global_stack.getProperty(key, "value"))

        return result

    ##  Get all extruder values for a certain setting. This function will skip the user settings container.
    #
    #   This is exposed to SettingFunction so it can be used in value functions.
    #
    #   \param key The key of the setting to retrieve values for.
    #
    #   \return A list of values for all extruders. If an extruder does not have a value, it will not be in the list.
    #           If no extruder has the value, the list will contain the global value.
    @staticmethod
    def getDefaultExtruderValues(key):
        global_stack = Application.getInstance().getGlobalContainerStack()
        context = PropertyEvaluationContext(global_stack)
        context.context["evaluate_from_container_index"] = 1  # skip the user settings container
        context.context["override_operators"] = {
            "extruderValue": ExtruderManager.getDefaultExtruderValue,
            "extruderValues": ExtruderManager.getDefaultExtruderValues,
            "resolveOrValue": ExtruderManager.getDefaultResolveOrValue
        }

        result = []
        for extruder in ExtruderManager.getInstance().getMachineExtruders(global_stack.getId()):
            # only include values from extruders that are "active" for the current machine instance
            if int(extruder.getMetaDataEntry("position")) >= global_stack.getProperty("machine_extruder_count", "value", context = context):
                continue

            value = extruder.getRawProperty(key, "value", context = context)

            if value is None:
                continue

            if isinstance(value, SettingFunction):
                value = value(extruder, context = context)

            result.append(value)

        if not result:
            result.append(global_stack.getProperty(key, "value", context = context))

        return result

    ##  Get all extruder values for a certain setting.
    #
    #   This is exposed to qml for display purposes
    #
    #   \param key The key of the setting to retrieve values for.
    #
    #   \return String representing the extruder values
    @pyqtSlot(str, result="QVariant")
    def getInstanceExtruderValues(self, key):
        return ExtruderManager.getExtruderValues(key)

    ##  Updates the material container to a material that matches the material diameter set for the printer
    def updateMaterialForDiameter(self, extruder_position: int, global_stack = None):
        if not global_stack:
            global_stack = Application.getInstance().getGlobalContainerStack()
            if not global_stack:
                return

        if not global_stack.getMetaDataEntry("has_materials", False):
            return

        extruder_stack = global_stack.extruders[str(extruder_position)]

        material_diameter = extruder_stack.material.getProperty("material_diameter", "value")
        if not material_diameter:
            # in case of "empty" material
            material_diameter = 0

        material_approximate_diameter = str(round(material_diameter))
        material_diameter = extruder_stack.definitionChanges.getProperty("material_diameter", "value")
        setting_provider = extruder_stack
        if not material_diameter:
            if extruder_stack.definition.hasProperty("material_diameter", "value"):
                material_diameter = extruder_stack.definition.getProperty("material_diameter", "value")
            else:
                material_diameter = global_stack.definition.getProperty("material_diameter", "value")
                setting_provider = global_stack

        if isinstance(material_diameter, SettingFunction):
            material_diameter = material_diameter(setting_provider)

        machine_approximate_diameter = str(round(material_diameter))

        if material_approximate_diameter != machine_approximate_diameter:
            Logger.log("i", "The the currently active material(s) do not match the diameter set for the printer. Finding alternatives.")

            if global_stack.getMetaDataEntry("has_machine_materials", False):
                materials_definition = global_stack.definition.getId()
                has_material_variants = global_stack.getMetaDataEntry("has_variants", False)
            else:
                materials_definition = "fdmprinter"
                has_material_variants = False

            old_material = extruder_stack.material
            search_criteria = {
                "type": "material",
                "approximate_diameter": machine_approximate_diameter,
                "material": old_material.getMetaDataEntry("material", "value"),
                "brand": old_material.getMetaDataEntry("brand", "value"),
                "supplier": old_material.getMetaDataEntry("supplier", "value"),
                "color_name": old_material.getMetaDataEntry("color_name", "value"),
                "definition": materials_definition
            }
            if has_material_variants:
                search_criteria["variant"] = extruder_stack.variant.getId()

            container_registry = Application.getInstance().getContainerRegistry()
            empty_material = container_registry.findInstanceContainers(id = "empty_material")[0]

            if old_material == empty_material:
                search_criteria.pop("material", None)
                search_criteria.pop("supplier", None)
                search_criteria.pop("brand", None)
                search_criteria.pop("definition", None)
                search_criteria["id"] = extruder_stack.getMetaDataEntry("preferred_material")

            materials = container_registry.findInstanceContainers(**search_criteria)
            if not materials:
                # Same material with new diameter is not found, search for generic version of the same material type
                search_criteria.pop("supplier", None)
                search_criteria.pop("brand", None)
                search_criteria["color_name"] = "Generic"
                materials = container_registry.findInstanceContainers(**search_criteria)
            if not materials:
                # Generic material with new diameter is not found, search for preferred material
                search_criteria.pop("color_name", None)
                search_criteria.pop("material", None)
                search_criteria["id"] = extruder_stack.getMetaDataEntry("preferred_material")
                materials = container_registry.findInstanceContainers(**search_criteria)
            if not materials:
                # Preferred material with new diameter is not found, search for any material
                search_criteria.pop("id", None)
                materials = container_registry.findInstanceContainers(**search_criteria)
            if not materials:
                # Just use empty material as a final fallback
                materials = [empty_material]

            Logger.log("i", "Selecting new material: %s", materials[0].getId())

            extruder_stack.material = materials[0]

    ##  Get the value for a setting from a specific extruder.
    #
    #   This is exposed to SettingFunction to use in value functions.
    #
    #   \param extruder_index The index of the extruder to get the value from.
    #   \param key The key of the setting to get the value of.
    #
    #   \return The value of the setting for the specified extruder or for the
    #   global stack if not found.
    @staticmethod
    def getExtruderValue(extruder_index, key):
        if extruder_index == -1:
            extruder_index = int(Application.getInstance().getMachineManager().defaultExtruderPosition)
        extruder = ExtruderManager.getInstance().getExtruderStack(extruder_index)

        if extruder:
            value = extruder.getRawProperty(key, "value")
            if isinstance(value, SettingFunction):
                value = value(extruder)
        else:
            # Just a value from global.
            value = Application.getInstance().getGlobalContainerStack().getProperty(key, "value")

        return value

    ##  Get the default value from the given extruder. This function will skip the user settings container.
    #
    #   This is exposed to SettingFunction to use in value functions.
    #
    #   \param extruder_index The index of the extruder to get the value from.
    #   \param key The key of the setting to get the value of.
    #
    #   \return The value of the setting for the specified extruder or for the
    #   global stack if not found.
    @staticmethod
    def getDefaultExtruderValue(extruder_index, key):
        extruder = ExtruderManager.getInstance().getExtruderStack(extruder_index)
        context = PropertyEvaluationContext(extruder)
        context.context["evaluate_from_container_index"] = 1  # skip the user settings container
        context.context["override_operators"] = {
            "extruderValue": ExtruderManager.getDefaultExtruderValue,
            "extruderValues": ExtruderManager.getDefaultExtruderValues,
            "resolveOrValue": ExtruderManager.getDefaultResolveOrValue
        }

        if extruder:
            value = extruder.getRawProperty(key, "value", context = context)
            if isinstance(value, SettingFunction):
                value = value(extruder, context = context)
        else:  # Just a value from global.
            value = Application.getInstance().getGlobalContainerStack().getProperty(key, "value", context = context)

        return value

    ##  Get the resolve value or value for a given key
    #
    #   This is the effective value for a given key, it is used for values in the global stack.
    #   This is exposed to SettingFunction to use in value functions.
    #   \param key The key of the setting to get the value of.
    #
    #   \return The effective value
    @staticmethod
    def getResolveOrValue(key):
        global_stack = Application.getInstance().getGlobalContainerStack()
        resolved_value = global_stack.getProperty(key, "value")

        return resolved_value

    ##  Get the resolve value or value for a given key without looking the first container (user container)
    #
    #   This is the effective value for a given key, it is used for values in the global stack.
    #   This is exposed to SettingFunction to use in value functions.
    #   \param key The key of the setting to get the value of.
    #
    #   \return The effective value
    @staticmethod
    def getDefaultResolveOrValue(key):
        global_stack = Application.getInstance().getGlobalContainerStack()
        context = PropertyEvaluationContext(global_stack)
        context.context["evaluate_from_container_index"] = 1  # skip the user settings container
        context.context["override_operators"] = {
            "extruderValue": ExtruderManager.getDefaultExtruderValue,
            "extruderValues": ExtruderManager.getDefaultExtruderValues,
            "resolveOrValue": ExtruderManager.getDefaultResolveOrValue
        }

        resolved_value = global_stack.getProperty(key, "value", context = context)

        return resolved_value
