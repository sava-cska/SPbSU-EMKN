package com.mcs.emkn.ui.courses.viewmodels

sealed class CoursesNavEvents {
    data class OpenCourse(val courseId: Int) : CoursesNavEvents()
}
