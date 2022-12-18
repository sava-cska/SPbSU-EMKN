package com.mcs.emkn.ui.changepassword

import android.app.AlertDialog
import android.os.Bundle
import android.os.CountDownTimer
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import com.mcs.emkn.R
import com.mcs.emkn.core.Router
import com.mcs.emkn.databinding.FragmentConfirmationBinding
import com.mcs.emkn.ui.changepassword.viewmodels.ChangePasswordError
import com.mcs.emkn.ui.changepassword.viewmodels.ChangePasswordInteractor
import com.mcs.emkn.ui.changepassword.viewmodels.ChangePasswordNavEvent
import com.mcs.emkn.ui.changepassword.viewmodels.ChangePasswordViewModel
import com.mcs.emkn.ui.emailconfirmation.viewmodels.EmailConfirmationError
import com.mcs.emkn.ui.emailconfirmation.viewmodels.EmailConfirmationNavEvent
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class ChangePasswordConfirmationFragment : Fragment() {
    private var _binding: FragmentConfirmationBinding? = null
    private val binding get() = _binding!!

    @Inject
    lateinit var router: Router
    private val changePasswordInteractor: ChangePasswordInteractor by viewModels<ChangePasswordViewModel>()

    private var verificationCode: String? = null
    private var timerStarted = false
    private var countDownTimer: CountDownTimer? = null

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = FragmentConfirmationBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        countDownTimer?.cancel()
        _binding = null
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        setupLayout()
        binding.sendCodeButton.setOnClickListener {
            verificationCode?.let { code -> changePasswordInteractor.validateCode(code) }
        }
        binding.sendCodeAgainButton.setOnClickListener {
            changePasswordInteractor.sendAnotherCode()
        }
        binding.backButton.setOnClickListener {
            onBackButtonPressed()
        }
        setupCodeEditField()

        subscribeToErrorsStatus()
        subscribeToNavStatus()
        subscribeToTimerStatus()
        loadTimer()
    }

    private fun setupLayout() {
        binding.confirmationHeader.text =
            resources.getString(R.string.change_password_confirmation_header)
    }

    private fun setupCodeEditField() {
        binding.sendCodeButton.isEnabled = false
        binding.codeEditText.onVerificationCodeChangeListener = { isFilled ->
            binding.sendCodeButton.isEnabled = isFilled
        }
        binding.codeEditText.onVerificationCodeFilledListener = { code ->
            verificationCode = code
        }
    }

    private fun onBackButtonPressed() {
        if (binding.codeEditText.isBlank()) {
            router.back()
            return
        }
        AlertDialog.Builder(requireContext())
            .setTitle(R.string.confirmation_back_alert_dialog_text)
            .setNegativeButton(R.string.sign_in_back_alert_dialog_cancel_button_text) { dialog, _ ->
                dialog?.dismiss()
            }
            .setPositiveButton(R.string.sign_in_back_alert_dialog_ok_button_text) { _, _ ->
                router.back()
            }
            .show()
    }

    private fun subscribeToErrorsStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                changePasswordInteractor.errors.collect { error ->
                    when (error) {
                        is ChangePasswordError.BadNetwork -> {
                            Toast
                                .makeText(
                                    requireContext(),
                                    resources.getString(R.string.bad_network_error),
                                    Toast.LENGTH_LONG
                                )
                                .show()
                        }
                        is ChangePasswordError.CodeExpired -> {
                            Toast
                                .makeText(
                                    requireContext(),
                                    resources.getString(R.string.code_expire_error),
                                    Toast.LENGTH_LONG
                                )
                                .show()
                        }
                        is ChangePasswordError.InvalidCode -> {
                            Toast
                                .makeText(
                                    requireContext(),
                                    resources.getString(R.string.code_invalid_error),
                                    Toast.LENGTH_LONG
                                )
                                .show()
                        }
                    }
                }
            }
        }
    }

    private fun subscribeToNavStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                changePasswordInteractor.navEvents.collect { navEvent ->
                    when (navEvent) {
                        is ChangePasswordNavEvent.ContinueChangePassword -> {
                            router.goToCommitChangePasswordScreen()
                        }
                    }
                }
            }
        }
    }

    private fun loadTimer() {
        lifecycleScope.launch {
            changePasswordInteractor.loadTimerAsync().await()?.let { timer ->
                timerStarted = true
                binding.sendCodeAgainButton.isVisible = false
                binding.timerTextVIew.isVisible = true
                countDownTimer?.cancel()
                startSendCodeTimer(timer)
            }
        }
    }


    private fun subscribeToTimerStatus() {
        lifecycleScope.launch {
            lifecycle.repeatOnLifecycle(Lifecycle.State.STARTED) {
                changePasswordInteractor.timerFlow.collect { timer ->
                    timerStarted = true
                    binding.sendCodeAgainButton.isVisible = false
                    binding.timerTextVIew.isVisible = true
                    countDownTimer?.cancel()
                    startSendCodeTimer(timer)
                }
            }
        }
    }

    private fun startSendCodeTimer(timeMills: Long) {
        countDownTimer = object : CountDownTimer(timeMills, 1000) {
            override fun onTick(millisUntilFinished: Long) {
                binding.timerTextVIew.text =
                    resources.getString(R.string.send_code_again_in, millisUntilFinished / 1000)
            }

            override fun onFinish() {
                timerStarted = false
                binding.sendCodeAgainButton.isVisible = true
                binding.timerTextVIew.isVisible = false
            }
        }.apply { start() }
    }
}