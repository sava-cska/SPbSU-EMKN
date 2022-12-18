package com.mcs.emkn.ui.coursehomeworks

import android.content.Context
import android.view.ViewGroup
import com.mcs.emkn.R
import com.mcs.emkn.core.rv.ViewHolder
import com.mcs.emkn.databinding.CourseHomeworksTaskViewHolderBinding

class CourseHomeworksTaskViewHolder(
    parent: ViewGroup,
) :
    ViewHolder<CourseHomeworksTaskItem>(R.layout.course_homeworks_task_view_holder, parent) {

    private val binding = CourseHomeworksTaskViewHolderBinding.bind(itemView)
    private val context: Context
        get() {
            return itemView.context
        }

    override fun bind(item: CourseHomeworksTaskItem) {
        binding.homeworkTitle.text = item.title
        binding.dateText.text = item.dayMonth
        binding.timeText.text = item.time
        when(item.status) {
            is CourseHomeworksTaskItem.HomeworkStatus.Checked -> {
                "${item.status.score}/${item.status.maximum}".also { binding.status.text = it }
            }
            is CourseHomeworksTaskItem.HomeworkStatus.NotChecked -> {
                binding.status.text = item.status.text
            }
            is CourseHomeworksTaskItem.HomeworkStatus.NotSubmitted -> {
                binding.status.text = item.status.text
            }
        }
    }
}