package com.mcs.emkn.ui.courses

import android.os.Bundle
import android.transition.TransitionManager
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.AdapterView
import android.widget.AdapterView.OnItemSelectedListener
import android.widget.ArrayAdapter
import androidx.activity.addCallback
import androidx.core.os.bundleOf
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import com.mcs.emkn.R
import com.mcs.emkn.core.Router
import com.mcs.emkn.core.rv.RecyclerAdapterWithDelegates
import com.mcs.emkn.core.rv.RouterBundleKeys
import com.mcs.emkn.core.rv.VerticalSpaceDecorator
import com.mcs.emkn.databinding.FragmentCoursesBinding
import com.mcs.emkn.ui.courses.viewmodels.CheckBoxesState
import com.mcs.emkn.ui.courses.viewmodels.Course
import com.mcs.emkn.ui.courses.viewmodels.CoursesInteractor
import com.mcs.emkn.ui.courses.viewmodels.CoursesViewModel
import com.mcs.emkn.ui.profile.viewmodels.Profile
import com.mcs.emkn.ui.profile.viewmodels.ProfilesLoader
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject


@AndroidEntryPoint
class CoursesFragment : Fragment(R.layout.fragment_courses) {
    private var _binding: FragmentCoursesBinding? = null
    private val binding get() = _binding!!
    private val adapter = RecyclerAdapterWithDelegates(
        listOf(CoursesAdapter(
            object : CoursesAdapter.CoursesListener {
                override fun onCourseEnrollButtonClick(courseId: Int) {
                    coursesInteractor.changeEnrollState(courseId)
                }

                override fun onCourseTitleClick(courseId: Int) {
                    val course = courses.find { it.id == courseId }
                    course?.let {
                        val title = course.title
                        val courseProfiles =
                            course.teachersProfiles.map { profile_id -> profiles[profile_id] }
                                .filterNotNull().toTypedArray()
                        val description = course.shortDescription
                        val bundle = bundleOf(
                            RouterBundleKeys.courseId to courseId,
                            RouterBundleKeys.title to title,
                            RouterBundleKeys.courseProfiles to courseProfiles,
                            RouterBundleKeys.description to description
                        )
                        router.goToCoursePage(bundle)
                    }
                }
            }
        )), listOf()
    )
    private val coursesInteractor: CoursesInteractor by viewModels<CoursesViewModel>()
    private var courses: List<Course> = listOf()
    private val profiles: MutableMap<Int, Profile> = mutableMapOf()
    private var checkBoxesState = CheckBoxesState(false, false)

    @Inject
    lateinit var profilesLoader: ProfilesLoader

    @Inject
    lateinit var router: Router

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        requireActivity().onBackPressedDispatcher.addCallback(this) {}
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = FragmentCoursesBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.coursesRecycler.adapter = adapter
        binding.coursesRecycler.addItemDecoration(
            VerticalSpaceDecorator(
                view.context.resources.getDimensionPixelSize(
                    R.dimen.courses_items_offset
                )
            )
        )

        lifecycleScope.launch {
            coursesInteractor.periods.collect {
                if (binding.coursesSettings.periodChooser.adapter?.isEmpty == false) {
                    return@collect
                }
                val periods = it.periods.reversed()
                val periodsNames = periods.map { it.text }
                val spinnerAdapter: ArrayAdapter<String> = ArrayAdapter(
                    this@CoursesFragment.requireContext(),
                    android.R.layout.simple_spinner_dropdown_item,
                    periodsNames
                )
                spinnerAdapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
                binding.coursesSettings.periodChooser.onItemSelectedListener =
                    object : OnItemSelectedListener {
                        override fun onItemSelected(
                            parent: AdapterView<*>?,
                            view: View?,
                            position: Int,
                            id: Long
                        ) {
                            coursesInteractor.onPeriodChosen(periods[position].id)
                        }

                        override fun onNothingSelected(parent: AdapterView<*>?) = Unit
                    }
                binding.coursesSettings.periodChooser.adapter = spinnerAdapter
            }
        }

        lifecycleScope.launch {
            coursesInteractor.courses.collect { it ->
                courses = it.data.courses
                binding.coursesLoader.visibility = if (it.hasMore) View.VISIBLE else View.GONE
                val profileIds = it.data.courses.mapNotNull { it.teachersProfiles.firstOrNull() }
                profilesLoader.requestProfiles(profileIds)
                updateList()
            }
        }
        lifecycleScope.launch {
            coursesInteractor.period.collect {
                binding.period.text = it.text
            }
        }
        lifecycleScope.launch {
            profilesLoader.profiles.collect {
                profiles.putAll(it)
                updateList()
            }
        }
        coursesInteractor.loadPeriodsAndCourses()
        binding.settingsIcon.setOnClickListener {
            TransitionManager.beginDelayedTransition(binding.coursesFragment)
            if (binding.coursesSettings.coursesSettings.visibility == View.VISIBLE) {
                binding.coursesSettings.coursesSettings.visibility = View.GONE
            } else {
                binding.coursesSettings.coursesSettings.visibility = View.VISIBLE
            }
        }
        binding.coursesSettings.excludeEnrolled.setOnCheckedChangeListener { buttonView, isChecked ->
            checkBoxesState = checkBoxesState.copy(isExcludingEnroll = isChecked)
            updateList()
        }
        binding.coursesSettings.excludeUnenrolled.setOnCheckedChangeListener { buttonView, isChecked ->
            checkBoxesState = checkBoxesState.copy(isExcludingUnenroll = isChecked)
            updateList()
        }
        binding.profileIcon.setOnClickListener {
            router.goToProfile()
        }
    }

    private fun updateList() {
        val courses = courses.filter {
            if (it.enrolled) {
                !checkBoxesState.isExcludingEnroll
            } else {
                !checkBoxesState.isExcludingUnenroll
            }
        }
        val items = courses.map {
            val teacher = it.teachersProfiles.firstOrNull().let { profiles[it] }
            val name = if (teacher == null) "..." else "${teacher.firstName} ${teacher.secondName}"
            CourseItem(
                it.id,
                it.title,
                name,
                it.shortDescription,
                if (it.enrolled) CourseItem.ButtonState.Unenroll else CourseItem.ButtonState.Enroll
            )
        }
        adapter.items = items
        adapter.notifyDataSetChanged()
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }

}