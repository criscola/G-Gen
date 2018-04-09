import QtQuick 2.2
import QtQuick.Controls 1.4
import QtQuick.Controls.Styles 1.4

import UM 1.2 as UM


Item
{
    id: extruderInfo
    property var printCoreConfiguration

    width: Math.round(parent.width / 2)
    height: childrenRect.height
    Label
    {
        id: materialLabel
        text: printCoreConfiguration.activeMaterial != null ? printCoreConfiguration.activeMaterial.name : ""
        elide: Text.ElideRight
        width: parent.width
        font: UM.Theme.getFont("very_small")
    }
    Label
    {
        id: printCoreLabel
        text: printCoreConfiguration.hotendID
        anchors.top: materialLabel.bottom
        elide: Text.ElideRight
        width: parent.width
        font: UM.Theme.getFont("very_small")
        opacity: 0.5
    }
}
