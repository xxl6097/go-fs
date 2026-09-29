#!/bin/bash
module=$(grep "module" go.mod | cut -d ' ' -f 2)
#options=("windows:amd64" "windows:arm64" "linux:amd64" "linux:arm64" "linux:arm:7" "linux:arm:5" "linux:mips64" "linux:mips64le" "linux:mips:softfloat" "linux:mipsle:softfloat" "linux:riscv64" "linux:loong64" "darwin:amd64" "darwin:arm64" "freebsd:amd64" "android:arm64")
options=("linux:arm64" "linux:amd64" "windows:amd64" "darwin:arm64")
version=$(git tag -l "v[0-99]*.[0-99]*.[0-99]*" --sort=-creatordate | head -n 1)
versionDir="$module/pkg"
#appname="tc-agent"
appname=$(basename $module)
isGarble="no"

function writeVersionGoFile() {
  if [ ! -d "./pkg" ]; then
    mkdir "./pkg"
  fi
cat <<EOF > ./pkg/version.go
package pkg

import (
	"fmt"
	"runtime"
	"strings"
)

func init() {
	OsType = runtime.GOOS
	Arch = runtime.GOARCH
	GoVersion = runtime.Version()
	Compiler = runtime.Compiler
}

var (
	AppName          string // 应用名称
	AppVersion       string // 应用版本
	BuildVersion     string // 编译版本
	BuildTime        string // 编译时间
	GoVersion        string // Golang信息
	DisplayName      string // 服务显示名
	Description      string // 服务描述信息
	OsType           string // 操作系统
	Arch             string // cpu类型
	Compiler         string // 编译器信息
	GitRevision      string // Git版本
	GitBranch        string // Git分支
	GitTreeState     string // state of git tree, either "clean" or "dirty"
	GitCommit        string // sha1 from git, output of 4a2ea0514582c5bdf629ad348341970c5ea8fdc6
	GitVersion       string // semantic version, derived by build scripts
	GitReleaseCommit string
	BinName          string // 运行文件名称，包含平台架构
	GithubUser       string // github用户
	GithubRepo       string // github项目名称
)

// Version 版本信息
func Version() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "App Name", AppName))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "App Version", AppVersion))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "DisplayName", DisplayName))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Description", Description))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Build version", BuildVersion))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Build time", BuildTime))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Golang Version", GoVersion))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "OsType", OsType))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Arch", Arch))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Compiler", Compiler))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Git revision", GitRevision))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "Git branch", GitBranch))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "GitTreeState", GitTreeState))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "GitCommit", GitCommit))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "GitVersion", GitVersion))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "GitReleaseCommit", GitReleaseCommit))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "BinName", BinName))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "GithubUser", GithubUser))
	sb.WriteString(fmt.Sprintf("%-16s: %-5s\n", "GithubRepo", GithubRepo))
	fmt.Println(sb.String())
	return sb.String()
}
EOF
}

# shellcheck disable=SC2120
function buildgo() {
  echo "=====$1,$2,$3,$4,$5,$6"
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  os=$5
  arch=$6
  extra=$7
  dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}
  flags='';
  if [ "${os}" = "linux" ] && [ "${arch}" = "arm" ] && [ "${extra}" != "" ] ; then
    if [ "${extra}" = "7" ]; then
      flags=GOARM=7;
      dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}hf
    elif [ "${extra}" = "5" ]; then
      flags=GOARM=5;
      dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}
    fi;
  elif [ "${os}" = "windows" ] ; then
    dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}.exe
    if [ "${arch}" = "amd64" ]; then
        go generate ${appdir}
    fi
  elif [ "${os}" = "linux" ] && ([ "${arch}" = "mips" ] || [ "${arch}" = "mipsle" ]) && [ "${extra}" != "" ] ; then
    flags=GOMIPS=${extra};
  fi;
  #echo "build：GOOS=${os} GOARCH=${arch} ${flags} ==> ${dstFilePath}"
  printf "build：GOOS=%-7s GOARCH=%-8s ==> %s\n" ${os} ${arch} ${dstFilePath}

  filename=$(basename "$dstFilePath")
  binName="-X '${versionDir}.BinName=${filename}'"
  #echo "--->env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} go build -trimpath -ldflags "$ldflags $binName -linkmode internal" -o ${dstFilePath} ${appdir}"
  #env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} go build -trimpath -ldflags "$ldflags $binName -linkmode internal" -o ${dstFilePath} ${appdir}
  if [ "$isGarble" = "yes" ]; then
    env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} garble -literals -tiny -seed=random build -trimpath -ldflags "$ldflags $binName -linkmode internal -buildid=" -buildvcs=false -o ${dstFilePath} ${appdir}
  else
    env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} go build -trimpath -ldflags "$ldflags $binName -linkmode internal" -o ${dstFilePath} ${appdir}
  fi
  #env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} garble -literals -tiny -seed=random build -trimpath -ldflags "$ldflags $binName -linkmode internal -buildid=" -buildvcs=false -o ${dstFilePath} ${appdir}
  if [ "${os}" = "windows" ] ; then
    if [ "${arch}" = "amd64" ]; then
        rm -rf ${appdir}/resource.syso
    fi
  fi;
}

# builddir：输出目录
# appname：应用名称
# version：应用版本
# appdir：main.go目录
# disname：显示名
# describe：描述
function buildMenu() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  disname=$5
  describe=$6
  ldflags=$(buildLdflags $appname $disname $describe)
  PS3="请选择需要编译的平台："
  select arch in "${options[@]}"; do
      if [[ -n "$arch" ]]; then
        IFS=":" read -r os arch extra <<< "$arch"
        buildgo $builddir $appname $version $appdir $os $arch $extra
        return $?
      else
        echo "输入无效，请重新选择。"
      fi
  done
}

# builddir：输出目录
# appname：应用名称
# version：应用版本
# appdir：main.go目录
# disname：显示名
# describe：描述
function buildAll() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  disname=$5
  describe=$6
  ldflags=$(buildLdflags $appname $disname $describe)
  for arch in "${options[@]}"; do
      IFS=":" read -r os arch extra <<< "$arch"
      buildgo $builddir $appname $version $appdir $os $arch $extra
  done
  #wait
}

version::get_version_vars() {
    # shellcheck disable=SC1083
    GIT_COMMIT="$(git rev-parse HEAD^{commit})"

    if git_status=$(git status --porcelain 2>/dev/null) && [[ -z ${git_status} ]]; then
        GIT_TREE_STATE="clean"
    else
        GIT_TREE_STATE="dirty"
    fi

    # stolen from k8s.io/hack/lib/version.sh
    # Use git describe to find the version based on annotated tags.
    if [[ -n ${GIT_VERSION-} ]] || GIT_VERSION=$(git describe --tags --abbrev=14 --match "v[0-9]*" "${GIT_COMMIT}" 2>/dev/null); then
        # This translates the "git describe" to an actual semver.org
        # compatible semantic version that looks something like this:
        #   v1.1.0-alpha.0.6+84c76d1142ea4d
        #
        # TODO: We continue calling this "git version" because so many
        # downstream consumers are expecting it there.
        # shellcheck disable=SC2001
        DASHES_IN_VERSION=$(echo "${GIT_VERSION}" | sed "s/[^-]//g")
        if [[ "${DASHES_IN_VERSION}" == "---" ]] ; then
            # We have distance to subversion (v1.1.0-subversion-1-gCommitHash)
            # shellcheck disable=SC2001
            GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-\([0-9]\{1,\}\)-g\([0-9a-f]\{14\}\)$/.\1\-\2/")
        elif [[ "${DASHES_IN_VERSION}" == "--" ]] ; then
            # We have distance to base tag (v1.1.0-1-gCommitHash)
            # shellcheck disable=SC2001
            GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-g\([0-9a-f]\{14\}\)$/-\1/")
        fi
        if [[ "${GIT_TREE_STATE}" == "dirty" ]]; then
            # git describe --dirty only considers changes to existing files, but
            # that is problematic since new untracked .go files affect the build,
            # so use our idea of "dirty" from git status instead.
            GIT_VERSION+="-dirty"
        fi


        # Try to match the "git describe" output to a regex to try to extract
        # the "major" and "minor" versions and whether this is the exact tagged
        # version or whether the tree is between two tagged versions.
        if [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?([-].*)?([+].*)?$ ]]; then
            GIT_MAJOR=${BASH_REMATCH[1]}
            GIT_MINOR=${BASH_REMATCH[2]}
            GIT_MINRR=${BASH_REMATCH[3]}
        fi

        # If GIT_VERSION is not a valid Semantic Version, then refuse to build.
        if ! [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
            echo "GIT_VERSION should be a valid Semantic Version. Current value: ${GIT_VERSION}"
            echo "Please see more details here: https://semver.org"
            exit 1
        fi
    fi

    GIT_RELEASE_TAG=$(git describe --abbrev=0 --tags)
    GIT_RELEASE_COMMIT=$(git rev-list -n 1  "${GIT_RELEASE_TAG}")
}

function buildLdflags() {
  local ldflags
  ldflags="-s -w"
  # shellcheck disable=SC2317
  function add_ldflag() {
      local key=${1}
      local val=${2}
      ldflags+=(
          "-X '${versionDir}.${key}=${val}'"
      )
  }
  #os_name=$(uname -s)
  #echo "os type $os_name"
  appname=$1
  DisplayName=$2
  Description=$3
  APP_NAME=${appname}
  BUILD_VERSION=$(if [ "$(git describe --tags --abbrev=0 2>/dev/null)" != "" ]; then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
  BUILD_TIME=$(TZ=Asia/Shanghai date "+%Y-%m-%d %H:%M:%S")
  GIT_REVISION=$(git rev-parse --short HEAD)
  GIT_BRANCH=$(git name-rev --name-only HEAD)
  #GIT_BRANCH=$(git tag -l "v[0-99]*.[0-99]*.[0-99]*" --sort=-creatordate | head -n 1)
  # shellcheck disable=SC2089
  version::get_version_vars
  add_ldflag "DisplayName" "${DisplayName}_${version}"
  add_ldflag "Description" "${Description}"
  add_ldflag "AppName" "${APP_NAME}"
  add_ldflag "AppVersion" "${version}"
  add_ldflag "BuildVersion" "${BUILD_VERSION}"
  add_ldflag "BuildTime" "${BUILD_TIME}"
  add_ldflag "GitRevision" "${GIT_REVISION}"
  add_ldflag "GitBranch" "${GIT_BRANCH}"
  add_ldflag "GitCommit" "${GIT_COMMIT}"
  add_ldflag "GitTreeState" "${GIT_TREE_STATE}"
  add_ldflag "GitVersion" "${GIT_VERSION}"
  add_ldflag "GitReleaseCommit" "${GIT_RELEASE_COMMIT}"
  echo "${ldflags[*]-}"
}


function showBuildDir() {
  # 检查是否输入路径参数
  if [ -z "$1" ]; then
      echo "用法: $0 <路径>"
      exit 1
  fi

  # 验证路径是否存在且为目录
  if [ ! -d "$1" ]; then
      echo "错误: 路径 '$1' 不存在或不是目录！"
      exit 1
  fi

  # 获取指定路径下的所有直接子目录（非递归）
  dirs=()
  while IFS= read -r dir; do
      dirs+=("$dir")
  done < <(find "$1" -maxdepth 1 -type d ! -path "$1" | sort)

  # 检查是否有子目录
  if [ ${#dirs[@]} -eq 0 ]; then
      echo "路径 '$1' 下没有子目录！"
      exit 0
  fi

  # 生成交互式菜单
  echo "请选择要操作的目录："
  PS3="输入序号 (1-${#dirs[@]}): "
  select dir in "${dirs[@]}"; do
      if [[ -n "$dir" ]] && [[ $REPLY -ge 1 && $REPLY -le ${#dirs[@]} ]]; then
          echo "您选择的目录是: $dir"
          break
#          return $dir
      else
          echo "无效输入！请输入有效序号"
      fi
  done
}


function install() {
 echo "${builddir}  ${appname}_${version}_${os}_${arch}"
 pwd
# host="v.uuxia.cn"
# host="192.168.10.7"
# host="uuxia.cn"
# host="10.6.14.26"
# host="192.168.1.3"
# bash <(curl -s -S -L http://${host}:8087/up) ${dstFilePath} /soft/${appname}/${version}
# sudo ${builddir}/${appname}_${version}_${os}_${arch} install
# ${builddir}/${appname}_${version}_${os}_${arch} install
}


# shellcheck disable=SC2120
function buildInstaller() {
#  version="v0.6.80"

  if [ -z "${version}" ]; then
      echo "---->${version}"
      version="v0.0.0"
  fi
  showBuildDir ./cmd
  builddir="./dist"
  #appname=$(basename "$dir")
#  appname="srvinstaller"
  appdir=${dir}
  disname="${appname}应用程序"
  describe="一款基于GO语言的服务安装程序"
  rm -rf ${builddir}
  buildMenu $builddir $appname "$version" $appdir $disname $describe
#  buildAll $builddir $appname "$version" $appdir $disname $describe
  #install
}

# shellcheck disable=SC2120
function githubActions() {
#  oldFile="./temp/sample.old"
#  newFile="./temp/sample.new"
#  patFile="./temp/sample.patch"
#  bsdiff ${oldFile} ${newFile} ${patFile}
#  ls -lh ./temp
  builddir="./release"
#  appname="srvinstaller"
  appdir="./cmd/app"
  disname="${appname}应用程序"
  describe="一款基于GO语言的服务安装程序"
  echo "===>version:${version}"
  go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
  go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
  rm -rf ${builddir}
  buildAll $builddir $appname "$version" $appdir $disname $describe
  mkdir -p ./release/packages
  mv -fv ./release/${appname}* ./release/packages
}

function showMenu() {
    echo "1. 编译程序"
    echo "2. 安装程序"
    echo "请输入编号:"
    read index
    clear
    case "$index" in
    [1]) (bootstrap $1 $2);;
    [2]) (install);;
    *) echo "exit" ;;
  esac
}

function check_garble() {
  # garble 必须在 PATH 中
  if ! command -v garble >/dev/null 2>&1; then
    echo "未找到 garble，正在安装..."
    go install mvdan.cc/garble@latest
    export PATH="$PATH:$(go env GOPATH)/bin"
  fi
}

function build_start() {
#  printf "是否garble混淆？(yes/no): "
#  read res
#  # 转小写
#  res_lower=$(echo "$res" | tr '[:upper:]' '[:lower:]')
#  if [ "$res_lower" = "yes" ]; then
#    isGarble="yes"
#    check_garble
#  fi
  buildInstaller
}

function bootstrap() {
  #printf "\033[1;31m%-10s\033[0m\n" "Error"  # 红色加粗文本
  if [ $# -ge 2 ] && [ -n "$2" ]; then
    version=$2
  fi
  writeVersionGoFile
  case $1 in
  github) (githubActions) ;;
    *) (build_start)  ;;
  esac
}
bootstrap $1 $2