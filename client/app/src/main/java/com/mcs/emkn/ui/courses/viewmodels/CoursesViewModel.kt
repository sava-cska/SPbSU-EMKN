package com.mcs.emkn.ui.courses.viewmodels

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.haroldadmin.cnradapter.NetworkResponse
import com.mcs.emkn.core.State
import com.mcs.emkn.database.Database
import com.mcs.emkn.database.entities.CheckBoxesStateEntity
import com.mcs.emkn.database.entities.CourseEntity
import com.mcs.emkn.database.entities.PeriodEntity
import com.mcs.emkn.database.entities.ProfileEntity
import com.mcs.emkn.network.Api
import com.mcs.emkn.network.dto.request.CoursesEnrollUnenrollRequestDto
import com.mcs.emkn.network.dto.request.CoursesListRequestDto
import com.mcs.emkn.network.dto.request.EmptyRequest
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableSharedFlow
import java.util.concurrent.Executors
import javax.inject.Inject

@HiltViewModel
class CoursesViewModel @Inject constructor(
        private val api: Api,
        private val db: Database,
) : CoursesInteractor, ViewModel() {
    override val courses: Flow<State<PeriodCourses>>
        get() = _coursesFlow
    override val periods: Flow<PeriodsData>
        get() = _periodsFlow
    override val navEventsFlow: Flow<CoursesNavEvents>
        get() = _navEvents
    override val errorsFlow: Flow<CoursesError>
        get() = _errorsFlow
    override val period: Flow<PeriodEntity>
        get() = _periodFlow

    private val dispatcher = Executors.newSingleThreadExecutor().asCoroutineDispatcher()
    private val scope = CoroutineScope(dispatcher)
    private var currentPeriod: Int? = null

    override fun onCleared() {
        scope.cancel()
        super.onCleared()
    }

    override fun onPeriodChosen(periodId: Int) {
        scope.launch {
            currentPeriod = periodId
            loadCourses()
            _periodFlow.emit(db.coursesDao().getPeriods().first { it.id == periodId })

        }
    }

    override fun loadPeriodsAndCourses() {
        scope.launch {
            loadCourses()
        }
    }

    override fun changeEnrollState(id: Int) {
        scope.launch {
            val course = db.coursesDao().getAllCourses().firstOrNull { it.id == id } ?: return@launch
            db.coursesDao().updateCourseEnrolledState(id, !course.enrolled)
            val auth = db.accountsDao().getCredentials().first().toAuthHeader()
            if (course.enrolled) {
                api.coursesUnenroll(CoursesEnrollUnenrollRequestDto(id), auth)
            } else {
                api.coursesEnroll(CoursesEnrollUnenrollRequestDto(id), auth)
            }
            loadCourses()
        }
    }

    private suspend fun loadCourses() {
        val local = currentPeriod
        val periodId = if (local == null) {
            loadAndMergePeriods()
            _periodFlow.emit(db.coursesDao().getPeriods().last())
            db.coursesDao().getPeriods().last().id
        } else {
            local
        }
        currentPeriod = periodId
        emitCourses(periodId, true)
        when (val response = api.coursesList(CoursesListRequestDto(listOf(periodId)), db.accountsDao().getCredentials().first().toAuthHeader())) {
            is NetworkResponse.Success -> {
                db.coursesDao().putCourses(response.body.response.coursesByPeriodDto.first().let { dto -> dto.courses.map { CourseEntity(dto.periodId, it.id, it.title, it.enrolled ?: false, it.shortDescription, it.teachersProfiles) } })
            }
            else -> Unit
        }
        emitCourses(periodId, false)
    }

    private suspend fun emitCourses(periodId: Int, hasMore: Boolean) {
        val newCourses = db.coursesDao().getCoursesByPeriods(listOf(periodId))
        val period = db.coursesDao().getPeriods().firstOrNull { it.id == periodId } ?: return
        _coursesFlow.emit(
                State(
                        PeriodCourses(
                                period, newCourses.map {
                            Course(it.id, it.title, it.enrolled, it.shortDescription, it.teachersProfiles)
                        }), hasMore
                )
        )
    }

    private suspend fun loadAndMergePeriods() {
        when (val response = api.coursesPeriods(EmptyRequest(), db.accountsDao().getCredentials().first().toAuthHeader())) {
            is NetworkResponse.Success -> {
                db.coursesDao().putPeriods(response.body.response.periods.map { PeriodEntity(it.id, it.text) })
            }
            else -> Unit
        }
        _periodsFlow.emit(PeriodsData(db.coursesDao().getPeriods()))
    }

    private val _coursesFlow = MutableSharedFlow<State<PeriodCourses>>()
    private val _periodsFlow = MutableSharedFlow<PeriodsData>()
    private val _navEvents = MutableSharedFlow<CoursesNavEvents>()
    private val _errorsFlow = MutableSharedFlow<CoursesError>()
    private val _periodFlow = MutableSharedFlow<PeriodEntity>()
}
