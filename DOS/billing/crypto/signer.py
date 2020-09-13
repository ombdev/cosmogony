import os
import random
import time
import subprocess
import string

from distutils.spawn import find_executable


class _HelperStr(object):

    @staticmethod
    def erase_bom(path):
        import codecs

        BUFSIZE = 4096
        chunk = None

        def takeout(l, f, c):
            i = 0
            c = c[l:]
            while c:
                f.seek(i)
                f.write(c)
                i += len(c)
                f.seek(l, os.SEEK_CUR)
                c = f.read(BUFSIZE)
            f.seek(-l, os.SEEK_CUR)
            f.truncate()

        with open(path, "r+b") as p:
            chunk = p.read(BUFSIZE)
            if chunk.startswith(codecs.BOM_UTF8):
                takeout(len(codecs.BOM_UTF8), p, chunk)
            if chunk.startswith(codecs.BOM_UTF32_BE):
                takeout(len(codecs.BOM_UTF32_BE), p, chunk)
            if chunk.startswith(codecs.BOM_UTF32_LE):
                takeout(len(codecs.BOM_UTF32_LE), p, chunk)
            if chunk.startswith(codecs.BOM_UTF16_BE):
                takeout(len(codecs.BOM_UTF16_BE), p, chunk)
            if chunk.startswith(codecs.BOM_UTF16_LE):
                takeout(len(codecs.BOM_UTF16_LE), p, chunk)

    @staticmethod
    def random_str(size=8):
        """generates random string as per size"""
        return ''.join(
            random.SystemRandom().choice(
                string.ascii_uppercase + string.digits
            ) for _ in range(size)
        )


class _LocalExec(object):

    def __init__(self, err_mute=False):
        self.__err_mute = err_mute

    def __call__(self, cmd_tokens, cmd_timeout, ign_rcs):
        """Execute a command on local machine."""

        def time_gap(delta):
            t = time.time()
            return t, t + delta

        def monitor(p, tbegin, tend):
            """Loop until process returns or timeout expires"""
            rc = None
            output = ''
            while time.time() < tend and rc is None:
                rc = p.poll()
                if rc is None:
                    try:
                        outs, errs = p.communicate(timeout=1)
                        output += outs
                    except subprocess.TimeoutExpired:
                        pass
            return output, rc

        if self.__err_mute:
            out_err = subprocess.DEVNULL
        else:
            out_err = subprocess.STDOUT

        output, rc = monitor(
            subprocess.Popen(
                cmd_tokens,
                universal_newlines=True,
                stdout=subprocess.PIPE,
                stderr= out_err
            ),
            *time_gap(cmd_timeout)
        )

        if rc is None:
            raise subprocess.TimeoutExpired(
                cmd=cmd_tokens,
                output=output,
                timeout=cmd_timeout
            )

        if ign_rcs is None:
            ign_rcs = []

        if rc in ign_rcs or rc == 0:
            return output

        raise subprocess.CalledProcessError(
            returncode=rc,
            cmd=cmd_tokens,
            output=output
        )


class SignerError(Exception):

    def __init__(self, message=None):
        self.message = message

    def __str__(self):
        return self.message


class Signer(object):
    """Callable object to sign string by openssl"""
    __SSL_BIN = "openssl"
    __SIZE_RANDOM_STR = 8
    __SUPPORTED = ['sha1', 'sha256']

    SHA1, SHA256 = range(2)

    def __init__(self, cipher, pem_privkey):

        # You must first extract the public key from the certificate:
        # openssl x509 -pubkey -noout -in cert.pem > pubkey.pem
        # then use the key to verify the signature:
        # openssl dgst -verify pubkey.pem -signature sigfile datafile

        def seekout_openssl():
            executable = find_executable(self.__SSL_BIN)
            if executable:
                return os.path.abspath(executable)
            raise SignerError("it has not found {} binary".format(self.__SSL_BIN))

        self._le = LocalExec()
        self._cipher = self.__SUPPORTED[cipher]
        self._pem_privkey = pem_privkey
        self._ssl_bin = seekout_openssl()

    def __call__(self, str2sign):
        """signs an string and returns base64 string"""
        tmp_dir = tempfile.gettempdir()
        sealbin_f = '{}/{}'.format(tmp_dir, _HelperStr.random_str(self.__SIZE_RANDOM_STR))
        input_f = '{}/{}'.format(tmp_dir, _HelperStr.random_str(self.__SIZE_RANDOM_STR))
        result_f = '{}/{}'.format(tmp_dir, _HelperStr.random_str(self.__SIZE_RANDOM_STR))

        self.__touch(input_f)

        with open(input_f, 'r+b') as cf:
            cf.write(str2sign.encode("utf-8-sig"))

        _HelperStr.erase_bom(input_f)

        dgst_args = [
            'dgst',
            '-{}'.format(self._cipher),
            '-sign',
            self._pem_privkey,
            '-out',
            sealbin_f,
            input_f
        ]

        base64_args = [
            'base64',
            '-in',
            sealbin_f,
            '-A',
            '-out',
            result_f
        ]

        try:
            self._le([self._ssl_bin] + dgst_args, cmd_timeout=10, ign_rcs=None)
            self._le([self._ssl_bin] + base64_args, cmd_timeout=10, ign_rcs=None)
        except subprocess.CalledProcessError as e:
            msg = "Command raised exception\nOutput: " + str(e.output)
            raise SignerError(msg)

        rs = self.__fetch_result(result_f)

        os.remove(sealbin_f)
        os.remove(input_f)
        os.remove(result_f)

        return rs

    def __touch(self, path):
        with open(path, 'a'):
            os.utime(path, None)

    def __fetch_result(self, path):
        rs = None
        statinfo = os.stat(path)
        if statinfo.st_size > 0:
            rs = ''
            with open(path, 'r') as rf:
                for line in rf:
                    rs = rs + line.replace("\n", "")
        if rs is None:
            SignerError("Unexpected ssl output!!!")
        return rs
