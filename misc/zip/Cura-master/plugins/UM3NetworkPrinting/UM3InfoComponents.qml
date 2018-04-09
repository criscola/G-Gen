import UM 1.2 as UM
import Cura 1.0 as Cura

import QtQuick 2.2
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.1
import QtQuick.Window 2.1

Item
{
    id: base

    property string activeQualityDefinitionId: Cura.MachineManager.activeQualityDefinitionId
    property bool isUM3: activeQualityDefinitionId == "ultimaker3" || activeQualityDefinitionId.match("ultimaker_") != null
    property bool printerConnected: Cura.MachineManager.printerConnected
    property bool printerAcceptsCommands: printerConnected && Cura.MachineManager.printerOutputDevices[0].acceptsCommands
    property bool authenticationRequested: printerConnected && (Cura.MachineManager.printerOutputDevices[0].authenticationState == 2 || Cura.MachineManager.printerOutputDevices[0].authenticationState == 5) // AuthState.AuthenticationRequested or AuthenticationReceived.

    Row
    {
        objectName: "networkPrinterConnectButton"
        visible: isUM3
        spacing: UM.Theme.getSize("default_margin").width

        Button
        {
            height: UM.Theme.getSize("save_button_save_to_button").height
            tooltip: catalog.i18nc("@info:tooltip", "Send access request to the printer")
            text: catalog.i18nc("@action:button", "Request Access")
            style: UM.Theme.styles.sidebar_action_button
            onClicked: Cura.MachineManager.printerOutputDevices[0].requestAuthentication()
            visible: printerConnected && !printerAcceptsCommands && !authenticationRequested
        }

        Button
        {
            height: UM.Theme.getSize("save_button_save_to_button").height
            tooltip: catalog.i18nc("@info:tooltip", "Connect to a printer")
            text: catalog.i18nc("@action:button", "Connect")
            style: UM.Theme.styles.sidebar_action_button
            onClicked: connectActionDialog.show()
            visible: !printerConnected
        }
    }

    UM.Dialog
    {
        id: connectActionDialog
        Loader
        {
            anchors.fill: parent
            source: "DiscoverUM3Action.qml"
        }
        rightButtons: Button
        {
            text: catalog.i18nc("@action:button", "Close")
            iconName: "dialog-close"
            onClicked: connectActionDialog.reject()
        }
    }


    Column
    {
        objectName: "networkPrinterConnectionInfo"
        visible: isUM3
        spacing: UM.Theme.getSize("default_margin").width
        anchors.fill: parent

        Button
        {
            tooltip: catalog.i18nc("@info:tooltip", "Send access request to the printer")
            text: catalog.i18nc("@action:button", "Request Access")
            onClicked: Cura.MachineManager.printerOutputDevices[0].requestAuthentication()
            visible: printerConnected && !printerAcceptsCommands && !authenticationRequested
        }

        Row
        {
            visible: printerConnected
            spacing: UM.Theme.getSize("default_margin").width

            anchors.left: parent.left
            anchors.right: parent.right
            height: childrenRect.height

            Column
            {
                Repeater
                {
                    model: Cura.ExtrudersModel { simpleNames: true }
                    Label { text: model.name }
                }
            }
            Column
            {
                Repeater
                {
                    id: nozzleColumn
                    model: printerConnected ? Cura.MachineManager.printerOutputDevices[0].hotendIds : null
                    Label { text: nozzleColumn.model[index] }
                }
            }
            Column
            {
                Repeater
                {
                    id: materialColumn
                    model: printerConnected ? Cura.MachineManager.printerOutputDevices[0].materialNames : null
                    Label { text: materialColumn.model[index] }
                }
            }
        }

        Button
        {
            tooltip: catalog.i18nc("@info:tooltip", "Load the configuration of the printer into Cura")
            text: catalog.i18nc("@action:button", "Activate Configuration")
            visible: false // printerConnected && !isClusterPrinter()
            onClicked: manager.loadConfigurationFromPrinter()
        }
    }

    UM.I18nCatalog{id: catalog; name:"cura"}
}
