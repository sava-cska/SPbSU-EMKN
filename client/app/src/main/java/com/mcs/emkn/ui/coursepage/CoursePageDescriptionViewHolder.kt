package com.mcs.emkn.ui.coursepage

import android.content.Context
import android.view.ViewGroup
import com.bumptech.glide.Glide
import com.bumptech.glide.load.resource.bitmap.RoundedCorners
import com.mcs.emkn.R
import com.mcs.emkn.core.rv.ViewHolder
import com.mcs.emkn.databinding.CoursePageDescriptionViewHolderBinding


class CoursePageDescriptionViewHolder(
    parent: ViewGroup,
) :
    ViewHolder<CoursePageDescriptionItem>(R.layout.course_page_description_view_holder, parent) {

    private val binding = CoursePageDescriptionViewHolderBinding.bind(itemView)

    override fun bind(item: CoursePageDescriptionItem) {
        binding.descriptionText.text = item.text
    }
}