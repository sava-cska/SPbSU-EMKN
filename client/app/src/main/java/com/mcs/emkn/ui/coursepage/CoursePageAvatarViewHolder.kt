package com.mcs.emkn.ui.coursepage

import android.content.Context
import android.view.RoundedCorner
import android.view.View
import android.view.ViewGroup
import androidx.appcompat.content.res.AppCompatResources
import com.bumptech.glide.Glide
import com.bumptech.glide.load.resource.bitmap.RoundedCorners
import com.mcs.emkn.R
import com.mcs.emkn.core.rv.ViewHolder
import com.mcs.emkn.databinding.CoursePageTeacherViewHolderBinding
import com.mcs.emkn.ui.courses.CourseItem
import com.mcs.emkn.ui.courses.CoursesAdapter


class CoursePageAvatarViewHolder(
    parent: ViewGroup,
) :
    ViewHolder<CoursePageAvatarItem>(R.layout.course_page_teacher_view_holder, parent) {

    private val binding = CoursePageTeacherViewHolderBinding.bind(itemView)
    private val context: Context
        get() {
            return itemView.context
        }

    override fun bind(item: CoursePageAvatarItem) {
        Glide.with(context).load(item.avatarUrl).centerCrop().transform(RoundedCorners(16)).error(R.drawable.person_icon)
            .into(binding.teacherAvatar)
        binding.teacherName.text = item.name
    }
}
