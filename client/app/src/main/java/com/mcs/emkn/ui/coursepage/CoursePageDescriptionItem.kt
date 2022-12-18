package com.mcs.emkn.ui.coursepage

import com.mcs.emkn.core.rv.Item

data class CoursePageDescriptionItem(
    val id: Int,
    val text: String
) : Item {
    override fun getItemId(): Long {
        return id.toLong()
    }
}