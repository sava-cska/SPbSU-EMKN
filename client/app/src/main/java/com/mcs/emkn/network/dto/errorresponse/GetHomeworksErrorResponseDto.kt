package com.mcs.emkn.network.dto.errorresponse

import com.mcs.emkn.network.dto.error.NetError
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class GetHomeworksErrorResponseDto (
    @Json(name = "invalid_course_id") val invalidCourseId: NetError?
)

