package com.mcs.emkn.network.dto.response

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class ResponseWithTokenAndTimeDto(
    @Json(name = "response") val tokenAndTimeDto: TokenAndTimeDto,
)

@JsonClass(generateAdapter = true)
data class TokenAndTimeDto(
    @Json(name = "expires_in") val expiresIn: String,
    @Json(name = "random_token") val randomToken: String,
)