package com.mcs.emkn.network.dto.response

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class CoursesPeriodsResponseDto(
    @Json(name = "response") val response: CoursesPeriods,
)
@JsonClass(generateAdapter = true)
data class CoursesPeriods(
    @Json(name = "periods") val periods: List<PeriodDto>,
)
@JsonClass(generateAdapter = true)
data class PeriodDto(
    @Json(name = "id") val id: Int,
    @Json(name = "text") val text: String
)
