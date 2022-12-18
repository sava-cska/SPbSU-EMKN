package com.mcs.emkn.ui.coursehomeworks.viewmodels

import com.mcs.emkn.core.State
import com.mcs.emkn.network.dto.response.CourseHomeworks
import com.mcs.emkn.ui.courses.viewmodels.PeriodCourses
import kotlinx.coroutines.flow.Flow

interface CourseHomeworksInteractor {
    val homeworksFlow: Flow<State<CourseHomeworks>>
    fun loadHomeworks(courseId: Int)
}