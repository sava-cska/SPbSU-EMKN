package com.mcs.emkn.network.dto.request

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class ValidateEmailRequestDto(
    @Json(name = "verification_code") val verification_code: String,
    @Json(name = "random_token") val random_token: String,
)