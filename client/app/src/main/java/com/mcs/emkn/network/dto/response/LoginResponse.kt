package com.mcs.emkn.network.dto.response

import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class LoginResponse(
    val id: Int,
)