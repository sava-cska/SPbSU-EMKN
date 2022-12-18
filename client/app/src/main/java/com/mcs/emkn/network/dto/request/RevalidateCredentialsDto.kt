package com.mcs.emkn.network.dto.request

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class RevalidateCredentialsDto(
    @Json(name = "random_token") val randomToken: String,
)