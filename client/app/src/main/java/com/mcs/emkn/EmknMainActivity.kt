package com.mcs.emkn

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.navigation.findNavController
import com.mcs.emkn.core.RouterImpl
import com.mcs.emkn.database.Database
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject


@AndroidEntryPoint
class EmknMainActivity : AppCompatActivity() {
    @Inject
    lateinit var routerImpl: RouterImpl

    @Inject
    lateinit var db: Database

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        routerImpl.attachNavController(findNavController(R.id.main_host_fragment))
        chooseNavGraph()
    }

    override fun onDestroy() {
        super.onDestroy()
        routerImpl.releaseNavController()
    }

    private fun chooseNavGraph() {
        routerImpl.goToRegistrationNavGraph()
    }
}