package com.mcs.emkn.ui.courses.viewmodels

import com.mcs.emkn.database.entities.PeriodEntity

data class PeriodCourses(
    val period: PeriodEntity,
    val courses: List<Course>
)