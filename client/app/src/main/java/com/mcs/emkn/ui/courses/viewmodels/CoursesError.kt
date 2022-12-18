package com.mcs.emkn.ui.courses.viewmodels

sealed interface CoursesError {
    object NetworkError : CoursesError
    class UnknownError(val msg: String) : CoursesError
    object CoursesListInvalidPeriodId : CoursesError
    object CoursesEnrollInvalidCourseId : CoursesError
    object CoursesEnrollAlreadyEnrolled : CoursesError
    object CoursesUnenrollInvalidCourseId : CoursesError
    object CoursesUnenrollCourseIsNotEnrolled : CoursesError
}