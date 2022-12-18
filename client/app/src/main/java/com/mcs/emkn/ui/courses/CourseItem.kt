package com.mcs.emkn.ui.courses

import com.mcs.emkn.core.rv.Item

data class CourseItem(
    val id: Int,
    val title: String,
    val lector: String,
    val description: String,
    val buttonState: ButtonState,
) : Item {
    sealed class ButtonState {
        object Enroll : ButtonState()
        object Unenroll : ButtonState()
        object Absence : ButtonState()
    }

    override fun getItemId(): Long {
        return id.toLong()
    }
}