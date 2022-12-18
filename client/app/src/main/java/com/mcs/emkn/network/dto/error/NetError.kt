package com.mcs.emkn.network.dto.error

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class NetError(
    @Json(name = "code") val code: String?,
)