package com.mcs.emkn.ui.coursehomeworks

import android.view.ViewGroup
import com.mcs.emkn.core.rv.Item
import com.mcs.emkn.core.rv.RecyclerDelegate


class CourseHomeworksSectionAdapter :
    RecyclerDelegate<CourseHomeworksSectionViewHolder, CourseHomeworksSectionItem> {
    override fun onBindViewHolder(viewHolder: CourseHomeworksSectionViewHolder, item: CourseHomeworksSectionItem) {
        viewHolder.bind(item)
    }

    override fun onCreateViewHolder(parent: ViewGroup): CourseHomeworksSectionViewHolder {
        return CourseHomeworksSectionViewHolder(parent)
    }

    override fun matchesItem(item: Item): Boolean {
        return item is CourseHomeworksSectionItem
    }
}