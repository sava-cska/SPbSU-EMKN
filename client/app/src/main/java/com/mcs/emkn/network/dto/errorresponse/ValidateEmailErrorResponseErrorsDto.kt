package com.mcs.emkn.network.dto.errorresponse

import com.mcs.emkn.network.dto.error.NetError
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass


@JsonClass(generateAdapter = true)
data class ValidateEmailErrorResponseErrorsDto(
    @Json(name = "code_invalid") val codeInvalid: NetError?,
    @Json(name = "registration_expired") val registrationExpired: NetError?,
)