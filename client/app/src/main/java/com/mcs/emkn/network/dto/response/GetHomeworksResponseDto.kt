package com.mcs.emkn.network.dto.response

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class GetHomeworksResponseDto(
    @Json(name = "response") val response: CourseHomeworks,
)

@JsonClass(generateAdapter = true)
data class CourseHomeworks(
    @Json(name = "homeworks") val homeworks: List<HomeworkDto>
)

@JsonClass(generateAdapter = true)
data class HomeworkDto(
    @Json(name = "id") val id: Int,
    @Json(name = "name") val name: String,
    @Json(name = "deadline") val deadline: Long,
    @Json(name = "status_not_passed") val statusNotPassed: TextStatus?,
    @Json(name = "status_unchecked") val statusUnchecked: TextStatus?,
    @Json(name = "status_checked") val statusChecked: ScoreStatus?,
)

@JsonClass(generateAdapter = true)
data class ScoreStatus(
    @Json(name = "total_score") val totalScore: Int,
    @Json(name = "score") val score: Int,
)

@JsonClass(generateAdapter = true)
data class TextStatus(
    @Json(name = "text") val text: String
)
