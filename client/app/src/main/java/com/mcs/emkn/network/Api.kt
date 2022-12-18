package com.mcs.emkn.network

import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.network.dto.errorresponse.*
import com.mcs.emkn.network.dto.request.*
import com.mcs.emkn.network.dto.response.*
import retrofit2.http.Body
import retrofit2.http.Header
import retrofit2.http.HeaderMap
import retrofit2.http.POST

interface Api {
    @POST("accounts/register")
    suspend fun accountsRegister(
        @Body request: RegistrationRequestDto
    ): NetworkResponse<ResponseWithTokenAndTimeDto, RegistrationErrorResponseDto>

    @POST("accounts/validate_email")
    suspend fun validateEmail(
        @Body request: ValidateEmailRequestDto
    ): NetworkResponse<Unit, ValidateEmailErrorResponseDto>

    @POST("accounts/login")
    suspend fun accountsLogin(
        @Body request: LoginRequestDto
    ): NetworkResponse<LoginResponse, LoginErrorResponseDto>

    @POST("accounts/begin_change_password")
    suspend fun accountsBeginChangePassword(
        @Body request: BeginChangePasswordRequestDto
    ): NetworkResponse<ResponseWithTokenAndTimeDto, BeginChangePasswordErrorResponseDto>

    @POST("accounts/validate_change_password")
    suspend fun accountsValidateChangePassword(
        @Body request: ValidateChangePasswordRequestDto
    ): NetworkResponse<ResponseWithTokenDto, ValidateChangePasswordErrorResponseDto>

    @POST("accounts/commit_change_password")
    suspend fun accountsCommitChangePassword(
        @Body request: CommitChangePasswordRequestDto
    ): NetworkResponse<Unit, CommitChangePasswordErrorResponseDto>

    @POST("accounts/revalidate_registration_credentials")
    suspend fun accountsRevalidateRegistrationCredentials(
        @Body request: RevalidateCredentialsDto,
    ): NetworkResponse<ResponseWithTokenAndTimeDto, RevalidateRegistrationCredentialsErrorResponseDto>

    @POST("accounts/revalidate_change_password_credentials")
    suspend fun accountsRevalidateChangePasswordCredentials(
        @Body request: RevalidateCredentialsDto,
    ): NetworkResponse<ResponseWithTokenAndTimeDto, RevalidateRegistrationCredentialsErrorResponseDto>

    @POST("courses/periods")
    suspend fun coursesPeriods(@Body empty: EmptyRequest, @Header("Authorization") auth: String) : NetworkResponse<CoursesPeriodsResponseDto, Unit>

    @POST("courses/list")
    suspend fun coursesList(
        @Body request: CoursesListRequestDto,
        @Header("Authorization") auth: String,
    ) : NetworkResponse<CoursesListResponseDto, CoursesListErrorResponseDto>

    @POST("courses/enroll")
    suspend fun coursesEnroll(
        @Body request: CoursesEnrollUnenrollRequestDto,
        @Header("Authorization") auth: String,
    ) : NetworkResponse<Unit, CoursesEnrollErrorResponseDto>

    @POST("courses/unenroll")
    suspend fun coursesUnenroll(
        @Body request: CoursesEnrollUnenrollRequestDto,
        @Header("Authorization") auth: String,
    ) : NetworkResponse<Unit, CoursesUnenrollErrorResponseDto>

    @POST("profiles/get")
    suspend fun profilesGet(
        @Body request: ProfilesGetRequestDto,
        @Header("Authorization") auth: String,
    ): NetworkResponse<ProfilesGetResponseDto, Unit>

    @POST("profiles/load_image")
    suspend fun uploadImage(
        @Body request: UploadImageRequestDto,
        @Header("Authorization") auth: String,
    ): NetworkResponse<Unit, Unit>

    @POST("courses/get_homeworks")
    suspend fun getHomeworks(
        @Body request: CoursesEnrollUnenrollRequestDto,
        @Header("Authorization") auth: String
    ): NetworkResponse<GetHomeworksResponseDto, GetHomeworksErrorResponseDto>
}
