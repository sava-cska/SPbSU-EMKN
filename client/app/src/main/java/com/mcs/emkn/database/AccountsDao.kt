package com.mcs.emkn.database

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.OnConflictStrategy.ABORT
import androidx.room.Query
import com.mcs.emkn.database.entities.ChangePasswordAttempt
import com.mcs.emkn.database.entities.ChangePasswordAttempt.Companion.CHANGE_PASSWORD_ATTEMPT_TABLE_NAME
import com.mcs.emkn.database.entities.ChangePasswordCommit
import com.mcs.emkn.database.entities.ChangePasswordCommit.Companion.CHANGE_PASSWORD_COMMIT_TABLE_NAME
import com.mcs.emkn.database.entities.Credentials
import com.mcs.emkn.database.entities.Credentials.Companion.CREDENTIALS_TABLE_NAME
import com.mcs.emkn.database.entities.SignUpAttempt
import com.mcs.emkn.database.entities.SignUpAttempt.Companion.SIGN_UP_ATTEMPT_TABLE_NAME

@Dao
interface AccountsDao {
    @Query("SELECT * FROM $CREDENTIALS_TABLE_NAME")
    fun getCredentials(): List<Credentials>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    fun putCredentials(credentials: Credentials)

    @Query("DELETE FROM $CREDENTIALS_TABLE_NAME")
    fun deleteCredentials()

    @Query("SELECT * FROM $SIGN_UP_ATTEMPT_TABLE_NAME")
    fun getSignUpAttempts(): List<SignUpAttempt>

    @Query("DELETE FROM $SIGN_UP_ATTEMPT_TABLE_NAME")
    fun deleteSignUpAttempts()

    @Insert(onConflict = ABORT)
    fun putSignUpAttempt(signUpAttempt: SignUpAttempt)

    @Insert(onConflict = ABORT)
    fun putChangePasswordAttempt(changePasswordAttempt: ChangePasswordAttempt)

    @Query("SELECT * FROM $CHANGE_PASSWORD_ATTEMPT_TABLE_NAME")
    fun getChangePasswordAttempts(): List<ChangePasswordAttempt>

    @Query("DELETE FROM $CHANGE_PASSWORD_ATTEMPT_TABLE_NAME")
    fun deleteChangePasswordAttempts()

    @Insert(onConflict = ABORT)
    fun putChangePasswordCommit(changePasswordCommit: ChangePasswordCommit)
    
    @Query("SELECT * FROM $CHANGE_PASSWORD_COMMIT_TABLE_NAME")
    fun getChangePasswordCommits(): List<ChangePasswordCommit>

    @Query("DELETE FROM $CHANGE_PASSWORD_COMMIT_TABLE_NAME")
    fun deleteChangePasswordCommits()
}