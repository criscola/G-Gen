# Copyright (c) 2017 Ultimaker B.V.
# Cura is released under the terms of the LGPLv3 or higher.

from collections import defaultdict
import threading
from typing import Any, Dict, Optional

from PyQt5.QtCore import pyqtProperty

from UM.Application import Application
from UM.Decorators import override

from UM.MimeTypeDatabase import MimeType, MimeTypeDatabase
from UM.Settings.ContainerStack import ContainerStack
from UM.Settings.SettingInstance import InstanceState
from UM.Settings.ContainerRegistry import ContainerRegistry
from UM.Settings.Interfaces import PropertyEvaluationContext
from UM.Logger import Logger

from . import Exceptions
from .CuraContainerStack import CuraContainerStack

##  Represents the Global or Machine stack and its related containers.
#
class GlobalStack(CuraContainerStack):
    def __init__(self, container_id: str, *args, **kwargs):
        super().__init__(container_id, *args, **kwargs)

        self.addMetaDataEntry("type", "machine")  # For backward compatibility

        self._extruders = {}  # type: Dict[str, "ExtruderStack"]

        # This property is used to track which settings we are calculating the "resolve" for
        # and if so, to bypass the resolve to prevent an infinite recursion that would occur
        # if the resolve function tried to access the same property it is a resolve for.
        # Per thread we have our own resolving_settings, or strange things sometimes occur.
        self._resolving_settings = defaultdict(set)  # keys are thread names

    ##  Get the list of extruders of this stack.
    #
    #   \return The extruders registered with this stack.
    @pyqtProperty("QVariantMap")
    def extruders(self) -> Dict[str, "ExtruderStack"]:
        return self._extruders

    @classmethod
    def getLoadingPriority(cls) -> int:
        return 2

    @classmethod
    def getConfigurationTypeFromSerialized(cls, serialized: str) -> Optional[str]:
        configuration_type = super().getConfigurationTypeFromSerialized(serialized)
        if configuration_type == "machine":
            return "machine_stack"
        return configuration_type

    ##  Add an extruder to the list of extruders of this stack.
    #
    #   \param extruder The extruder to add.
    #
    #   \throws Exceptions.TooManyExtrudersError Raised when trying to add an extruder while we
    #                                            already have the maximum number of extruders.
    def addExtruder(self, extruder: ContainerStack) -> None:
        position = extruder.getMetaDataEntry("position")
        if position is None:
            Logger.log("w", "No position defined for extruder {extruder}, cannot add it to stack {stack}", extruder = extruder.id, stack = self.id)
            return

        if any(item.getId() == extruder.id for item in self._extruders.values()):
            Logger.log("w", "Extruder [%s] has already been added to this stack [%s]", extruder.id, self.getId())
            return

        self._extruders[position] = extruder
        Logger.log("i", "Extruder[%s] added to [%s] at position [%s]", extruder.id, self.id, position)

    ##  Overridden from ContainerStack
    #
    #   This will return the value of the specified property for the specified setting,
    #   unless the property is "value" and that setting has a "resolve" function set.
    #   When a resolve is set, it will instead try and execute the resolve first and
    #   then fall back to the normal "value" property.
    #
    #   \param key The setting key to get the property of.
    #   \param property_name The property to get the value of.
    #
    #   \return The value of the property for the specified setting, or None if not found.
    @override(ContainerStack)
    def getProperty(self, key: str, property_name: str, context: Optional[PropertyEvaluationContext] = None) -> Any:
        if not self.definition.findDefinitions(key = key):
            return None

        if context is None:
            context = PropertyEvaluationContext()
        context.pushContainer(self)

        # Handle the "resolve" property.
        if self._shouldResolve(key, property_name, context):
            current_thread = threading.current_thread()
            self._resolving_settings[current_thread.name].add(key)
            resolve = super().getProperty(key, "resolve", context)
            self._resolving_settings[current_thread.name].remove(key)
            if resolve is not None:
                return resolve

        # Handle the "limit_to_extruder" property.
        limit_to_extruder = super().getProperty(key, "limit_to_extruder", context)
        if limit_to_extruder is not None:
            if limit_to_extruder == -1:
                limit_to_extruder = int(Application.getInstance().getMachineManager().defaultExtruderPosition)
            limit_to_extruder = str(limit_to_extruder)
        if limit_to_extruder is not None and limit_to_extruder != "-1" and limit_to_extruder in self._extruders:
            if super().getProperty(key, "settable_per_extruder", context):
                result = self._extruders[str(limit_to_extruder)].getProperty(key, property_name, context)
                if result is not None:
                    context.popContainer()
                    return result
            else:
                Logger.log("e", "Setting {setting} has limit_to_extruder but is not settable per extruder!", setting = key)

        result = super().getProperty(key, property_name, context)
        context.popContainer()
        return result

    ##  Overridden from ContainerStack
    #
    #   This will simply raise an exception since the Global stack cannot have a next stack.
    @override(ContainerStack)
    def setNextStack(self, next_stack: ContainerStack) -> None:
        raise Exceptions.InvalidOperationError("Global stack cannot have a next stack!")

    # protected:

    # Determine whether or not we should try to get the "resolve" property instead of the
    # requested property.
    def _shouldResolve(self, key: str, property_name: str, context: Optional[PropertyEvaluationContext] = None) -> bool:
        if property_name is not "value":
            # Do not try to resolve anything but the "value" property
            return False

        current_thread = threading.current_thread()
        if key in self._resolving_settings[current_thread.name]:
            # To prevent infinite recursion, if getProperty is called with the same key as
            # we are already trying to resolve, we should not try to resolve again. Since
            # this can happen multiple times when trying to resolve a value, we need to
            # track all settings that are being resolved.
            return False

        setting_state = super().getProperty(key, "state", context = context)
        if setting_state is not None and setting_state != InstanceState.Default:
            # When the user has explicitly set a value, we should ignore any resolve and
            # just return that value.
            return False

        return True


## private:
global_stack_mime = MimeType(
    name = "application/x-cura-globalstack",
    comment = "Cura Global Stack",
    suffixes = ["global.cfg"]
)

MimeTypeDatabase.addMimeType(global_stack_mime)
ContainerRegistry.addContainerTypeByName(GlobalStack, "global_stack", global_stack_mime.name)
