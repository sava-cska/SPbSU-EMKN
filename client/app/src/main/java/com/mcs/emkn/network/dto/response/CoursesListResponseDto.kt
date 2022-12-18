package com.mcs.emkn.network.dto.response

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class CoursesListResponseDto(
    @Json(name = "response") val response: CoursesList,
)

@JsonClass(generateAdapter = true)
data class CoursesList(
    @Json(name = "courses_by_period") val coursesByPeriodDto: List<PeriodCoursesDto>
)

@JsonClass(generateAdapter = true)
data class PeriodCoursesDto(
    @Json(name = "period_id") val periodId: Int,
    @Json(name = "courses") val courses: List<CourseDto>
)

@JsonClass(generateAdapter = true)
data class CourseDto(
    @Json(name = "id") val id: Int,
    @Json(name = "title") val title: String,
    @Json(name = "enrolled") val enrolled: Boolean?,
    @Json(name = "short_description") val shortDescription: String,
    @Json(name = "teachers_profiles") val teachersProfiles: List<Int>
)

