package com.mcs.emkn.network.dto.errorresponse

import com.mcs.emkn.network.dto.error.NetError
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class BeginChangePasswordErrorResponseErrorsDto(
    @Json(name = "illegal_email") val illegalEmail: NetError?,
)

