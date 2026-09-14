package app.everyreview.presentation.common

import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import app.everyreview.BuildConfig
import app.everyreview.R

/** Human-readable version string shown in the UI, e.g. `v0.1.0-dev.be4efb5 (4)`. */
val appVersionLabel: String
    get() = "v${BuildConfig.VERSION_NAME} (${BuildConfig.VERSION_CODE})"

@Composable
fun AppLogo(modifier: Modifier = Modifier, size: Int = 48) {
    Image(
        painter = painterResource(R.drawable.ic_logo),
        contentDescription = stringResource(R.string.app_name),
        modifier = modifier.size(size.dp),
    )
}

/** Logo + app name + version, for placing over the scan preview or on an about screen. */
@Composable
fun BrandHeader(modifier: Modifier = Modifier, onDark: Boolean = false) {
    val textColor = if (onDark) Color.White else MaterialTheme.colorScheme.onSurface
    Row(
        modifier = modifier,
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp),
    ) {
        AppLogo(size = 40)
        Column {
            Text(
                text = stringResource(R.string.app_name),
                style = MaterialTheme.typography.titleMedium,
                color = textColor,
            )
            Text(
                text = appVersionLabel,
                style = MaterialTheme.typography.labelSmall,
                color = textColor.copy(alpha = 0.75f),
            )
        }
    }
}

/** [BrandHeader] on a translucent pill so it stays readable on top of the camera preview. */
@Composable
fun BrandOverlay(modifier: Modifier = Modifier) {
    BrandHeader(
        onDark = true,
        modifier = modifier
            .background(Color.Black.copy(alpha = 0.45f), MaterialTheme.shapes.large)
            .padding(horizontal = 12.dp, vertical = 8.dp),
    )
}
