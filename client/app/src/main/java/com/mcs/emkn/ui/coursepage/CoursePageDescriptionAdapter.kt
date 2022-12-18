package com.mcs.emkn.ui.coursepage

import android.view.ViewGroup
import com.mcs.emkn.core.rv.Item
import com.mcs.emkn.core.rv.RecyclerDelegate



class CoursePageDescriptionAdapter() :
    RecyclerDelegate<CoursePageDescriptionViewHolder, CoursePageDescriptionItem> {
    override fun onBindViewHolder(viewHolder: CoursePageDescriptionViewHolder, item: CoursePageDescriptionItem) {
        viewHolder.bind(item)
    }

    override fun onCreateViewHolder(parent: ViewGroup): CoursePageDescriptionViewHolder {
        return CoursePageDescriptionViewHolder(parent)
    }

    override fun matchesItem(item: Item): Boolean {
        return item is CoursePageDescriptionItem
    }
}