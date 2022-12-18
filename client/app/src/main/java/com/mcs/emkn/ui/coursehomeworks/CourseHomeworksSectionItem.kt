package com.mcs.emkn.ui.coursehomeworks

import com.mcs.emkn.core.rv.Item

data class CourseHomeworksSectionItem(
    val id: Int,
    val title: String
) : Item {
    override fun getItemId(): Long {
        return id.toLong()
    }
}