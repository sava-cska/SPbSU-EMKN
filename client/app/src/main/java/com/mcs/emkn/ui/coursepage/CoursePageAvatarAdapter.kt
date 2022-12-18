package com.mcs.emkn.ui.coursepage

import android.view.ViewGroup
import com.mcs.emkn.core.rv.Item
import com.mcs.emkn.core.rv.RecyclerDelegate
import com.mcs.emkn.ui.courses.CourseItem
import com.mcs.emkn.ui.courses.CourseViewHolder


class CoursePageAvatarAdapter() : RecyclerDelegate<CoursePageAvatarViewHolder, CoursePageAvatarItem> {
    override fun onBindViewHolder(viewHolder: CoursePageAvatarViewHolder, item: CoursePageAvatarItem) {
        viewHolder.bind(item)
    }

    override fun onCreateViewHolder(parent: ViewGroup): CoursePageAvatarViewHolder {
        return CoursePageAvatarViewHolder(parent)
    }

    override fun matchesItem(item: Item): Boolean {
        return item is CoursePageAvatarItem
    }
}