package com.mcs.emkn.ui.coursehomeworks

import android.content.Context
import android.view.ViewGroup
import com.mcs.emkn.R
import com.mcs.emkn.core.rv.ViewHolder
import com.mcs.emkn.databinding.CourseHomeworksSectionViewHolderBinding
import com.mcs.emkn.databinding.CourseHomeworksTaskViewHolderBinding



class CourseHomeworksSectionViewHolder(
    parent: ViewGroup,
) :
    ViewHolder<CourseHomeworksSectionItem>(R.layout.course_homeworks_section_view_holder, parent) {

    private val binding = CourseHomeworksSectionViewHolderBinding.bind(itemView)
    private val context: Context
        get() {
            return itemView.context
        }

    override fun bind(item: CourseHomeworksSectionItem) {
        binding.sectionName.text = item.title
    }
}