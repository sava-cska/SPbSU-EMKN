package com.mcs.emkn.ui.courses

import android.content.Context
import android.view.View
import android.view.ViewGroup
import androidx.appcompat.content.res.AppCompatResources
import com.mcs.emkn.R
import com.mcs.emkn.core.rv.ViewHolder
import com.mcs.emkn.databinding.CourseViewHolderBinding

class CourseViewHolder(
        parent: ViewGroup,
        private val listener: CoursesAdapter.CoursesListener,
) :
    ViewHolder<CourseItem>(R.layout.course_view_holder, parent) {

    private val binding = CourseViewHolderBinding.bind(itemView)
    private val context: Context
        get() {
            return itemView.context
        }

    override fun bind(item: CourseItem) {
        binding.courseName.text = item.title
        binding.lector.text = item.lector
        binding.description.text = item.description
        when (item.buttonState) {
            CourseItem.ButtonState.Absence -> binding.enrollButton.visibility = View.GONE
            CourseItem.ButtonState.Enroll -> {
                binding.enrollButton.visibility = View.VISIBLE
                binding.enrollButton.text = context.getString(R.string.enroll_button)
                binding.enrollButton.setTextColor(context.getColor(R.color.text_reversed))
                binding.enrollButton.background =
                    AppCompatResources.getDrawable(context, R.drawable.confirm_background)
            }
            CourseItem.ButtonState.Unenroll -> {
                binding.enrollButton.visibility = View.VISIBLE
                binding.enrollButton.text = context.getString(R.string.unenroll_button)
                binding.enrollButton.setTextColor(context.getColor(R.color.confirm))
                binding.enrollButton.background =
                    AppCompatResources.getDrawable(context, R.drawable.unenroll_background)
            }
        }
        binding.enrollButton.setOnClickListener {
            listener.onCourseEnrollButtonClick(item.id)
        }
        binding.courseName.setOnClickListener {
            listener.onCourseTitleClick(item.id)
        }
    }
}