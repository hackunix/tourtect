package com.tourtect.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.tourtect.feature.assistant.AssistantRoute
import dagger.hilt.android.AndroidEntryPoint

private data class PrimaryDestination(val route: String, val label: String)

private val primaryDestinations = listOf(
    PrimaryDestination("assistant", "Assistant"),
    PrimaryDestination("explore", "Explore"),
    PrimaryDestination("community", "Community"),
    PrimaryDestination("saved", "Saved"),
    PrimaryDestination("profile", "Profile")
)

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme {
                Surface(modifier = Modifier.fillMaxSize()) {
                    TourtectApp()
                }
            }
        }
    }
}

@Composable
private fun TourtectApp() {
    val navController = rememberNavController()
    val backStackEntry by navController.currentBackStackEntryAsState()
    val currentRoute = backStackEntry?.destination?.route

    Scaffold(
        topBar = { GlobalActionBar(navController) },
        bottomBar = {
            NavigationBar {
                primaryDestinations.forEach { destination ->
                    NavigationBarItem(
                        selected = currentRoute == destination.route,
                        onClick = { navController.navigatePrimary(destination.route) },
                        icon = { Text(destination.label.take(1)) },
                        label = { Text(destination.label) }
                    )
                }
            }
        }
    ) { padding ->
        NavHost(
            navController = navController,
            startDestination = "assistant",
            modifier = Modifier.padding(padding)
        ) {
            composable("assistant") {
                AssistantRoute(onOpenDestination = { navController.navigateKnownTarget(it) })
            }
            composable("explore") {
                UnavailableScreen(
                    title = "Explore",
                    message = "Place exploration is not connected in this Android build."
                )
            }
            composable("community") {
                UnavailableScreen(
                    title = "Community knowledge",
                    message = "Community remains a secondary knowledge surface; Android feed wiring is pending."
                )
            }
            composable("saved") {
                UnavailableScreen("Saved", "Saved items are not connected in this Android build.")
            }
            composable("profile") {
                UnavailableScreen("Profile", "Profile and secure account storage are not connected yet.")
            }
            composable("live") {
                UnavailableScreen(
                    title = "Live voice unavailable",
                    message = "Microphone capture and realtime assistant-session handoff are not enabled. No audio is being recorded."
                )
            }
            composable("lens") {
                UnavailableScreen(
                    title = "Lens unavailable",
                    message = "A consent-bound backend capture contract is not enabled. No image has been captured or uploaded."
                )
            }
            composable("sos") {
                UnavailableScreen(
                    title = "Emergency directory unavailable",
                    message = "This build will not invent emergency numbers or place a call automatically."
                )
            }
            composable("price-check") {
                UnavailableScreen(
                    title = "Manual Price Check",
                    message = "The deterministic Price Check form is not connected in this Android build."
                )
            }
            composable("safety") {
                UnavailableScreen(
                    title = "Manual Safety Assessment",
                    message = "The rule-first Safety Assessment form is not connected in this Android build."
                )
            }
        }
    }
}

@Composable
private fun GlobalActionBar(navController: NavHostController) {
    Row(
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically,
        modifier = Modifier.fillMaxWidth().padding(horizontal = 12.dp, vertical = 6.dp)
    ) {
        Text("Tourtect", style = MaterialTheme.typography.titleLarge)
        Row {
            TextButton(onClick = { navController.navigateKnownTarget("live") }) { Text("Live") }
            TextButton(onClick = { navController.navigateKnownTarget("lens") }) { Text("Lens") }
            Button(onClick = { navController.navigateKnownTarget("sos") }) { Text("SOS") }
        }
    }
}

@Composable
private fun UnavailableScreen(title: String, message: String) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
        modifier = Modifier.fillMaxSize().padding(24.dp)
    ) {
        Text(title, style = MaterialTheme.typography.headlineSmall)
        Text(message, modifier = Modifier.padding(top = 8.dp))
    }
}

private fun NavHostController.navigatePrimary(route: String) {
    navigate(route) {
        popUpTo("assistant") { saveState = true }
        launchSingleTop = true
        restoreState = true
    }
}

private fun NavHostController.navigateKnownTarget(target: String) {
    val normalized = when {
        target.contains("price", ignoreCase = true) -> "price-check"
        target.contains("safety", ignoreCase = true) -> "safety"
        target.contains("community", ignoreCase = true) || target.contains("posts", ignoreCase = true) -> "community"
        target.contains("place", ignoreCase = true) || target.contains("explore", ignoreCase = true) -> "explore"
        target.equals("live", ignoreCase = true) -> "live"
        target.equals("lens", ignoreCase = true) -> "lens"
        target.contains("sos", ignoreCase = true) || target.contains("emergency", ignoreCase = true) -> "sos"
        else -> "assistant"
    }
    navigate(normalized) { launchSingleTop = true }
}
