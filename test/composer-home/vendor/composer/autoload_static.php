<?php

// autoload_static.php @generated by Composer

namespace Composer\Autoload;

class ComposerStaticInit8420ffedece46256b03cac1adb1a7f6f
{
    public static $prefixLengthsPsr4 = array (
        'H' => 
        array (
            'Hirak\\Prestissimo\\' => 18,
        ),
        'D' => 
        array (
            'Damoon\\Cito\\' => 12,
        ),
    );

    public static $prefixDirsPsr4 = array (
        'Hirak\\Prestissimo\\' => 
        array (
            0 => __DIR__ . '/..' . '/hirak/prestissimo/src',
        ),
        'Damoon\\Cito\\' => 
        array (
            0 => __DIR__ . '/..' . '/damoon/cito/src',
        ),
    );

    public static function getInitializer(ClassLoader $loader)
    {
        return \Closure::bind(function () use ($loader) {
            $loader->prefixLengthsPsr4 = ComposerStaticInit8420ffedece46256b03cac1adb1a7f6f::$prefixLengthsPsr4;
            $loader->prefixDirsPsr4 = ComposerStaticInit8420ffedece46256b03cac1adb1a7f6f::$prefixDirsPsr4;

        }, null, ClassLoader::class);
    }
}
