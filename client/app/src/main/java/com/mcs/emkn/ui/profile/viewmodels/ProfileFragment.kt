package com.mcs.emkn.ui.profile.viewmodels

import android.app.Activity.RESULT_OK
import android.content.Intent
import android.graphics.Bitmap
import android.net.Uri
import android.os.Bundle
import android.provider.MediaStore
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import com.bumptech.glide.Glide
import com.mcs.emkn.core.Router
import com.mcs.emkn.databinding.ProfileFragmentBinding
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject


@AndroidEntryPoint
class ProfileFragment : Fragment() {
    @Inject
    lateinit var router: Router

    private val binding: ProfileFragmentBinding
        get() {
            return _binding!!
        }
    private var _binding: ProfileFragmentBinding? = null

    private val profileViewModel by viewModels<ProfileViewModel>()

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = ProfileFragmentBinding.inflate(inflater)
        return binding.profileFragment
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        lifecycleScope.launch {
            profileViewModel.profile.collect {
                if (it.avatarUri.toString().isEmpty()) {
                    return@collect
                }
                Glide.with(this@ProfileFragment.requireContext()).load(it.avatarUri)
                    .into(binding.avatar)
            }
        }
        lifecycleScope.launch {
            profileViewModel.loadProfile()
        }
        binding.enrollButton.setOnClickListener {
            val intent = Intent(Intent.ACTION_GET_CONTENT)
            intent.type = "image/*"
            startActivityForResult(intent, PICK_CODE)

        }
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCode, resultCode, data)
        if (requestCode != PICK_CODE) {
            return
        }
        if (resultCode != RESULT_OK) {
            return
        }
        if (data == null ) {
            return
        }
        val selectedImage: Uri = data.data ?: return
        val bitmap: Bitmap = MediaStore.Images.Media.getBitmap(this.requireActivity().contentResolver, selectedImage)
        profileViewModel.uploadPhoto(bitmap)
    }

    companion object {
        private const val PICK_CODE = 30
    }
}