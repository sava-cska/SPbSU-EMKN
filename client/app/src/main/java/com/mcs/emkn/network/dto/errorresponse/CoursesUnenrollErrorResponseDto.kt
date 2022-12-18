package com.mcs.emkn.network.dto.errorresponse

import com.mcs.emkn.network.dto.error.NetError
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class CoursesUnenrollErrorResponseDto (
    @Json(name = "course_is_not_enrolled") val courseIsNotEnrolled: NetError?,
    @Json(name = "invalid_course_id") val invalidCourseId: NetError?
)
