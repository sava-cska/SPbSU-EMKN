package com.mcs.emkn.ui.coursehomeworks

import android.view.ViewGroup
import com.mcs.emkn.core.rv.Item
import com.mcs.emkn.core.rv.RecyclerDelegate


class CourseHomeworksTaskAdapter :
    RecyclerDelegate<CourseHomeworksTaskViewHolder, CourseHomeworksTaskItem> {
    override fun onBindViewHolder(viewHolder: CourseHomeworksTaskViewHolder, item: CourseHomeworksTaskItem) {
        viewHolder.bind(item)
    }

    override fun onCreateViewHolder(parent: ViewGroup): CourseHomeworksTaskViewHolder {
        return CourseHomeworksTaskViewHolder(parent)
    }

    override fun matchesItem(item: Item): Boolean {
        return item is CourseHomeworksTaskItem
    }
}