package com.mcs.emkn.ui.profile.viewmodels

import android.graphics.Bitmap
import android.util.Base64
import androidx.lifecycle.ViewModel
import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.database.Database
import com.mcs.emkn.network.Api
import com.mcs.emkn.network.dto.request.UploadImageRequestDto
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.newSingleThreadContext
import java.io.ByteArrayOutputStream
import javax.inject.Inject


@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val db: Database,
    private val profilesLoader: ProfilesLoader,
    private val api: Api,
) : ViewModel() {
    @OptIn(DelicateCoroutinesApi::class)
    private val scope = CoroutineScope(newSingleThreadContext("profile_worker"))

    val profile: Flow<Profile>
        get() = _profile

    private val _profile: MutableSharedFlow<Profile> = MutableSharedFlow()
    private val id by lazy {
        db.accountsDao().getCredentials().first().id !!
    }

    init {
        scope.launch {
            profilesLoader.profiles.collect {
                val profile = it[id] ?: return@collect
                _profile.emit(profile)
            }
        }
    }

    fun loadProfile() {
        scope.launch {
            profilesLoader.requestProfiles(listOf(id))
        }
    }

    fun uploadPhoto(bitmap: Bitmap) {
        scope.launch {
            val byteArrayOutputStream = ByteArrayOutputStream()
            bitmap.compress(Bitmap.CompressFormat.JPEG, 100, byteArrayOutputStream)
            val byteArray: ByteArray = byteArrayOutputStream.toByteArray()
            val encoded: String = Base64.encodeToString(byteArray, Base64.DEFAULT)
            when (api.uploadImage(UploadImageRequestDto(encoded), db.accountsDao().getCredentials().first().toAuthHeader())) {
                is NetworkResponse.Success -> {
                    loadProfile()
                }
                else -> {}
            }
        }
    }
}