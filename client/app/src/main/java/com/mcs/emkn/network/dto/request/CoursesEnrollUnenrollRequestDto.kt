package com.mcs.emkn.network.dto.request

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class CoursesEnrollUnenrollRequestDto (
    @Json(name = "course_id") val courseId: Int
)
