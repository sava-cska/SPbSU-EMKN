package com.mcs.emkn.ui.profile.viewmodels

import android.net.Uri
import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.database.Database
import com.mcs.emkn.database.entities.ProfileEntity
import com.mcs.emkn.network.Api
import com.mcs.emkn.network.dto.request.ProfilesGetRequestDto
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.update
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ProfilesLoader @Inject constructor(
    private val db: Database,
    private val api: Api,
) {

    private val _profiles = MutableStateFlow<Map<Int, Profile>>(mapOf())

    val profiles: Flow<Map<Int, Profile>>
        get() = _profiles

    @OptIn(ExperimentalCoroutinesApi::class)
    private val scope = CoroutineScope(Dispatchers.IO.limitedParallelism(1))

    fun requestProfiles(ids: List<Int>) {
        scope.launch {
            val profiles = db.coursesDao().getProfilesByIds(ids).toMutableList()
            _profiles.emit(profiles.associate { it.id to Profile(it.id, Uri.parse(it.avatarUrl), it.firstName, it.secondName) })
            when (val response = api.profilesGet(ProfilesGetRequestDto(ids), db.accountsDao().getCredentials().first().toAuthHeader())) {
                is NetworkResponse.Success -> {
                    db.coursesDao().putProfiles(response.body.response.profiles.map { ProfileEntity(it.id, it.avatarUrl, it.firstName, it.secondName) })
                    val profiles = db.coursesDao().getProfilesByIds(ids).toMutableList()
                    _profiles.emit(profiles.associate { it.id to Profile(it.id, Uri.parse(it.avatarUrl), it.firstName, it.secondName) })
                }
                else -> {}
            }
        }
    }
}