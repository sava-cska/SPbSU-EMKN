package com.mcs.emkn.network.dto.request

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class CoursesListRequestDto (
    @Json(name = "period_ids") val periodIds: List<Int>
)