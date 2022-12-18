package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.SignUpAttempt.Companion.SIGN_UP_ATTEMPT_TABLE_NAME

@Entity(tableName = SIGN_UP_ATTEMPT_TABLE_NAME)
data class SignUpAttempt(
    val email: String,
    @PrimaryKey
    val login: String,
    val password: String,
    val name: String,
    val surName: String,
    val randomToken: String,
    val expiresInSeconds: Long,
    val createdAt: Long,
) {
    companion object {
        const val SIGN_UP_ATTEMPT_TABLE_NAME = "sign_up_attempt"
    }
}