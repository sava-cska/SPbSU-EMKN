package com.mcs.emkn.ui.courses

import android.view.ViewGroup
import com.mcs.emkn.core.rv.Item
import com.mcs.emkn.core.rv.RecyclerDelegate

class CoursesAdapter(
        private val listener: CoursesListener
) : RecyclerDelegate<CourseViewHolder, CourseItem> {

    interface CoursesListener {
        fun onCourseEnrollButtonClick(courseId: Int)

        fun onCourseTitleClick(courseId: Int)
    }

    override fun onBindViewHolder(viewHolder: CourseViewHolder, item: CourseItem) {
        viewHolder.bind(item)
    }

    override fun onCreateViewHolder(parent: ViewGroup): CourseViewHolder {
        return CourseViewHolder(parent, listener)
    }

    override fun matchesItem(item: Item): Boolean {
        return item is CourseItem
    }
}