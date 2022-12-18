package com.mcs.emkn.network.dto.response

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class ResponseWithTokenDto(
    @Json(name = "response") val token: TokenDto,
)

@JsonClass(generateAdapter = true)
data class TokenDto(
    @Json(name = "change_password_token") val changePasswordToken: String,
)