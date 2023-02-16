#!/usr/bin/env perl

use File::Temp;
use Socket;

my $tmpfile = mktemp("$ENV{TMPDIR}tmpXXXXXX");
system("go", "build", "-o", "$tmpfile", ".");

my $path = "test.sock";

if (-S $path) {
    unlink($path);
}

$^F = 3;
socket(my $sock, PF_UNIX, SOCK_STREAM, 0) || die "Could not create socket: $!";
bind($sock, sockaddr_un($path)) || die "Could not bind socket: $!";
listen($sock, SOMAXCONN) || die "Could not listen on socket: $!";

$ENV{"LISTEN_PID"} = "$$";
$ENV{"LISTEN_FDS"} = "1";

exec("$tmpfile", "activate");
