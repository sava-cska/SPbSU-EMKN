package com.mcs.emkn.network.dto.request

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class CommitChangePasswordRequestDto(
    @Json(name = "change_password_token") val changePasswordToken: String,
    @Json(name = "new_password") val newPassword: String,
)