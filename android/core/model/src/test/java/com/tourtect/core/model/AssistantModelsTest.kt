package com.tourtect.core.model

import org.junit.Assert.assertEquals
import org.junit.Test

class AssistantModelsTest {
    @Test
    fun `wire values map to constrained enums`() {
        assertEquals(AssistantIntent.PRICE_CHECK, AssistantIntent.fromWire("price_check"))
        assertEquals(AssistantSafetyState.NON_EMERGENCY, AssistantSafetyState.fromWire("non_emergency"))
    }

    @Test
    fun `unknown values abstain instead of selecting a capability`() {
        assertEquals(AssistantIntent.UNKNOWN, AssistantIntent.fromWire("arbitrary_tool"))
        assertEquals(AssistantSafetyState.UNKNOWN, AssistantSafetyState.fromWire("certainly_safe"))
    }
}
