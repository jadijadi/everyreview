package app.everyreview.data.media

import android.content.Context
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.graphics.Matrix
import android.media.ExifInterface
import android.net.Uri
import androidx.core.content.FileProvider
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.io.File
import java.io.IOException
import java.util.UUID
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Turns a camera/gallery [Uri] into a small JPEG on disk that fits the backend's 5 MB cap
 * (and doesn't burn the user's data plan): longest edge capped, EXIF rotation baked in.
 */
@Singleton
class PhotoPreparer @Inject constructor(
    @ApplicationContext private val context: Context,
) {
    private val cameraDir get() = File(context.cacheDir, "camera").apply { mkdirs() }
    private val uploadDir get() = File(context.cacheDir, "uploads").apply { mkdirs() }

    /** A writable target for `ActivityResultContracts.TakePicture`, shared via the manifest's FileProvider. */
    fun newCaptureUri(): Uri {
        val file = File(cameraDir, "capture-${UUID.randomUUID()}.jpg")
        return FileProvider.getUriForFile(context, "${context.packageName}.fileprovider", file)
    }

    suspend fun prepareForUpload(source: Uri): File = withContext(Dispatchers.IO) {
        val bounds = BitmapFactory.Options().apply { inJustDecodeBounds = true }
        openStream(source).use { BitmapFactory.decodeStream(it, null, bounds) }
        if (bounds.outWidth <= 0 || bounds.outHeight <= 0) throw IOException("Could not read the selected image")

        val options = BitmapFactory.Options().apply { inSampleSize = sampleSize(bounds.outWidth, bounds.outHeight) }
        val decoded = openStream(source).use { BitmapFactory.decodeStream(it, null, options) }
            ?: throw IOException("Could not decode the selected image")
        val upright = applyExifRotation(source, decoded)

        val target = File(uploadDir, "upload-${UUID.randomUUID()}.jpg")
        target.outputStream().use { upright.compress(Bitmap.CompressFormat.JPEG, JPEG_QUALITY, it) }
        if (upright !== decoded) decoded.recycle()
        upright.recycle()
        target
    }

    private fun openStream(uri: Uri) =
        context.contentResolver.openInputStream(uri) ?: throw IOException("Could not open the selected image")

    private fun sampleSize(width: Int, height: Int): Int {
        var sample = 1
        while (maxOf(width, height) / (sample * 2) >= MAX_EDGE_PX) sample *= 2
        return sample
    }

    private fun applyExifRotation(uri: Uri, bitmap: Bitmap): Bitmap {
        val orientation = runCatching {
            openStream(uri).use { ExifInterface(it).getAttributeInt(ExifInterface.TAG_ORIENTATION, ExifInterface.ORIENTATION_NORMAL) }
        }.getOrDefault(ExifInterface.ORIENTATION_NORMAL)
        val degrees = when (orientation) {
            ExifInterface.ORIENTATION_ROTATE_90 -> 90f
            ExifInterface.ORIENTATION_ROTATE_180 -> 180f
            ExifInterface.ORIENTATION_ROTATE_270 -> 270f
            else -> return bitmap
        }
        val matrix = Matrix().apply { postRotate(degrees) }
        return Bitmap.createBitmap(bitmap, 0, 0, bitmap.width, bitmap.height, matrix, true)
    }

    private companion object {
        const val MAX_EDGE_PX = 1600
        const val JPEG_QUALITY = 85
    }
}
