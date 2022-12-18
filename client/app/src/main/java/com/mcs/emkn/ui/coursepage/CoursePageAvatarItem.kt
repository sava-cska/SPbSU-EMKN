package com.mcs.emkn.ui.coursepage

import com.mcs.emkn.core.rv.Item

data class CoursePageAvatarItem(
    val id: Int,
    val avatarUrl: String,
    val name: String
) : Item {
    override fun getItemId(): Long {
        return id.toLong()
    }
}